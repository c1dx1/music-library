package services

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"music-library/internal/models"
	"music-library/internal/repositories"
	pagination "music-library/internal/utils"
)

type MLibService struct {
	repo         repositories.MLibRepository
	log          *logrus.Logger
	extAPIClient *ExternalAPIClient
}

func NewMLibService(repo repositories.MLibRepository, log *logrus.Logger, extAPIClient *ExternalAPIClient) *MLibService {
	return &MLibService{repo: repo, log: log, extAPIClient: extAPIClient}
}

func (s *MLibService) GetLibrary(filter models.LibraryFilter) ([]models.Song, error) {
	s.log.Debug("Entering MLibService.GetLibrary func")

	if filter.Page != nil {
		if *filter.Page < 1 {
			*filter.Page = 1
			s.log.Info("Page is less than 1, defaulting to 1")
		}
	} else {
		page := 1
		filter.Page = &page
		s.log.Info("Page is nil, defaulting to 1")
	}

	if filter.Limit != nil {
		if *filter.Limit < 1 {
			*filter.Limit = 10
			s.log.Info("Limit is less than 1, defaulting to 10")
		}
	} else {
		limit := 10
		filter.Limit = &limit
		s.log.Info("Limit is nil, defaulting to 10")
	}

	filters := make(map[string]interface{})

	if filter.ID != nil {
		filters["s.id"] = filter.ID
	}
	if filter.Group != nil {
		filters["g.group_name"] = "%" + *filter.Group + "%"
	}
	if filter.Song != nil {
		filters["s.song_name"] = "%" + *filter.Song + "%"
	}
	if filter.ReleaseDate != nil {
		filters["s.release_date"] = *filter.ReleaseDate
	}
	if filter.Text != nil {
		filters["s.text"] = "%" + *filter.Text + "%"
	}
	if filter.Link != nil {
		filters["s.link"] = *filter.Link
	}

	songs, err := s.repo.GetLibrary(context.Background(), filters, *filter.Page, *filter.Limit)
	if err != nil {
		s.log.Debug("MLibService.GetLibrary err")
		return nil, fmt.Errorf("mlib service: getLib: repo: ", err)
	}
	s.log.Debug("MLibService.GetLibrary success")

	return songs, nil
}

func (s *MLibService) GetText(id, page, limit int) ([]string, error) {
	s.log.Debug("Entering MLibService.GetText func")

	text, err := s.repo.GetText(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("mlib service: getText: repo: ", err)
	}

	verses, err := pagination.PaginateVerses(text, page, limit, s.log)
	if err != nil {
		return nil, fmt.Errorf("mlib service: getText: paginateVerses: repo: ", err)
	}

	s.log.Debug("MLibService.GetText success")
	return verses, nil
}

func (s *MLibService) DeleteSong(id int) error {
	s.log.Debug("Entering MLibService.DeleteSong func")

	err := s.repo.DeleteSong(context.Background(), id)
	if err != nil {
		return fmt.Errorf("mlib service: DeleteSong: repo: ", err)
	}

	s.log.Debug("MLibService.DeleteSong success")

	return nil
}

func (s *MLibService) EditSong(id int, req models.EditSong) error {
	updates := make(map[string]interface{})

	if req.Group != nil {
		updates["group_name"] = *req.Group
	}
	if req.Song != nil {
		updates["song_name"] = *req.Song
	}
	if req.ReleaseDate != nil {
		updates["release_date"] = *req.ReleaseDate
	}
	if req.Text != nil {
		updates["text"] = *req.Text
	}
	if req.Link != nil {
		updates["link"] = *req.Link
	}

	if len(updates) == 0 {
		s.log.Infof("EditSong called with no updates for song ID: %d", id)
		return nil
	}

	s.log.Debugf("EditSong: Preparing to update song ID %d with updates: %+v", id, updates)

	err := s.repo.EditSong(context.Background(), id, updates)
	if err != nil {
		s.log.Errorf("EditSong: Failed to update song ID %d: %v", id, err)
		return fmt.Errorf("mlib service: EditSong: repo: %w", err)
	}

	s.log.Infof("EditSong: Successfully updated song ID %d", id)
	return nil
}

func (s *MLibService) AddSong(song models.Song) error {
	s.log.Infof("AddSong: Adding new song: %+v", song)

	err := s.extAPIClient.GetSongDetails(&song)
	if err != nil {
		s.log.Errorf("AddSong: Failed to get song details for: %+v, error: %v", song, err)
		return fmt.Errorf("mlib service: AddSong: GetSongDetails: %w", err)
	}

	s.log.Debugf("AddSong: Retrieved song details: %+v", song)

	err = s.repo.AddSong(context.Background(), song)
	if err != nil {
		s.log.Errorf("AddSong: Failed to add song to repository: %+v, error: %v", song, err)
		return fmt.Errorf("mlib_service: AddSong: repo: %w", err)
	}

	s.log.Infof("AddSong: Successfully added new song: %+v", song)
	return nil
}
