INSERT INTO groups (group_name)
VALUES ('The Beatles'),
       ('Queen'),
       ('Nirvana'),
       ('Led Zeppelin');

INSERT INTO songs (group_id, song_name, release_date, text, link)
VALUES (1, 'Hey Jude', '1968-08-26', 'Hey, Jude, don''t make it bad\nTake a sad song and make it better\nRemember to let her into your heart\nThen you can start to make it better\n\nHey, Jude, don''t be afraid You were made to go out and get her\nThe minute you let her under your skin\nThen you begin to make it better\nAnd anytime you feel the pain, hey, Jude, refrain\nDon''t carry the world upon your shoulders\nFor well you know that it''s a fool who plays it cool\nBy making his world a little colder\nNa-na-na-na-na, na-na-na-na\n\nHey, Jude, don''t let me down\n\nYou have found her, now go and get her\n(Let it out and let it in)\nRemember (Hey, Jude) to let her into your heart\nThen you can start to make it better',
        'https://example.com/heyjude'),
       (1, 'Let It Be', '1970-03-06', 'When I find myself in times of trouble\nMother Mary comes to me\nSpeaking words of wisdom, "Let it be".\nAnd in my hour of darkness\nShe is standing right in front of me\nSpeaking words of wisdom, "Let it be".\n\nLet it be, let it be,\nLet it be, let it be.\nWhisper words of wisdom, let it be.\n\nAnd when the broken-hearted people\nLiving in the world agree\nThere will be an answer, let it be.\nFor though they may be parted,\nThere is still a chance that they will see.\nThere will be an answer, let it be.\n\nLet it be, let it be,\nLet it be, let it be.\nThere will be an answer, let it be.',
        'https://example.com/letitbe'),
       (2, 'Bohemian Rhapsody', '1975-10-31', 'Is this the real life?\nIs this just fantasy?Caught in a landslide\nNo escape from reality\nOpen your eyes\nLook up to the skies and see\nI''m just a poor boy, I need no sympathy\nBecause I''m easy come, easy go\nA little high, little low\nAnyway the wind blows, doesn''t really matter to me, to me\n\nMama, just killed a man\nPut a gun against his head\nPulled my trigger, now he''s dead\nMama, life had just begun\nBut now I''ve gone and thrown it all away\nMama, ooo\nDidn''t mean to make you cry\nIf I''m not back again this time tomorrow\nCarry on, carry on, as if nothing really matters\n\nToo late, my time has come\nSends shivers down my spine\nBody''s aching all the time\nGoodbye everybody - I''ve got to go\nGotta leave you all behind and face the truth\nMama, ooo - (anyway the wind blows)\nI don''t want to die\nI sometimes wish I''d never been born at all',
        'https://example.com/bohemianrhapsody'),
       (3, 'Smells Like Teen Spirit', '1991-09-10', 'Load up on guns and bring your friends\nIt''s fun to lose and to pretend\nShe''s over-bored and self-assured\nOh no, I know a dirty word\n\nHello, hello, hello, how low?\nHello, hello, hello, how low?\nHello, hello, hello, how low?\nHello, hello, hello\n\nWith the lights out, it''s less dangerous\nHere we are now, entertain us\nI feel stupid and contagious\nHere we are now, entertain us',
        'https://example.com/smellsliketeenspirit'),
       (4, 'Stairway to Heaven', '1971-11-08',
        'There''s a lady who''s sure\nall that glitters is gold,\nand she''s buying a stairway to heaven\nWhen she gets there she knows,\nif the stores are all closed,\nwith a word she can get what\nshe came for.\n\nOoh, ooh,\nand she''s buying a stairway to heaven.\n\nThere''s a sign on the wall\nbut she wants to be sure,\n''cause you know sometimes\nwords have two meanings.\nIn a tree by the brook there''s a songbird who sings,\nSometimes all of our thoughts are misgiven.',
        'https://example.com/stairwaytoheaven');