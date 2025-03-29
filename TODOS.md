- [x] Invite link doesnt work directly if user is registering using that link.

  - using session to check for invite code.
  - is not retained when redirect to game page is used after register.
  - can use local storage to retain the invite code. every time user loads the page, check local storage for invite code. if present, call invite friend api, make the friend and flush the invite code from localstorage.

- [x] Register user - add text to use authentic email since email validation by sending confirmation mail wil be done.

- [ ] Login after playing game:

  - game should start without even asking for login or register. player can play the game, and then to save their score, can login in/register.
  - should think of securing the score saving in client side. So that user cannot manipulate the score from browser application side (inspect element).

- [x] Fix created_at timestamp update to null issue while answering questions.

- [ ] Some way to store which questions have already been attempted by the player, and when getting set of 5 question at random, add exception clause for the attempted questions. If less than 5 non attempted questions are left, send all of them. if none are left. Send error with "quiz completed by the user".

- [x] Background image change with fade throttle.

- [x] Save logs to file

- [x] Bug: The options generated from 15 random destinations (excuding the 5 questions) are not random. They are the first non included 15 questions from the table, which is repetitive. ORDER BY RANDOM() does not work in supabase go client.

- [x] Migrate to raw postgres connection rather than using supabase client for go.
