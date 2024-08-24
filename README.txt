Notes on using this program:

Open the directory containing docker-compose.yml and Dockerfile inside of command line after pulling this from GitHub
Simply type "docker-compose up --build" and hit enter
This should run the program and it should tell you if it is running or if it exits with a 1.
If it is running, open a browser tab and go to "http://localhost:8000"
You should see a sign in page. Simply type in an email and password, hit sign up, and receive a message at the bottom if successful or if there is a problem
Click the "already have an account" button to switch to the log in side of the page.
Retype the user email and password you just created in the step before this last one
Hit 'Log In" button
If password or email is invalid, it will not proceed and tell you that the password or email is incorrect
If it is correct, then it moves on to a new page with text saying "Congrats!" and a log out button
If you hit the logout button, then you will return to the sign in page




KNOWN ISSUES:
Sometimes the database is fails to connect. I had gone to sleep after getting it working only to wake up and find it was no longer working. Editing the connection string fixed it.
A user can hit the back arrow on a browser after logout and just return to the congrats screen without having to log in again.
The database user requires a password, this could lead to a failure to connect, though I have the password baked into the docker-compose file and the main.go. This is most likely not ideal but I could not find a way around it.
The Log in/Sign up page are the same page, simply put together with a toggle function. Sign up will always show first due to how this is set up. Most likely not ideal as well as you ALWAYS have to manually click "Already have an account" to then log in.