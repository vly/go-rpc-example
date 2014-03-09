# Assignment 1: RPC

## General notes
- Objective is to gain an understanding of using RPC.
- During week 1's lecture, the head tutor confirmed to the class that you can use whichever language you are comfortable with. Furthermore, confirmed I can use Golang when quested specifically.
- Have 4 weeks to complete this project

## Submission output
- a clear and brief (less than one page) description of what I did, and how the program works
- a record of running the application 9e.g. produced by the script command)
- source code and any corresponding make file
- altogether submitted in a single zipfile via weblearn

## Task 1
- Create a simlified version of a live message service. It delivers text messages from users to other users, provided both sender and receiver are connected to the system at that point in time. Message delivery is according to the pull paradigm, that is, messages have to be fetched by users; the service does not send a notification about the arrival of a message.
- End-user can activate the following operations through a command interface (e.g. IRC's /connect etc.):
	- connect to the system
	- disconnect
	- deposit a message for another connected user
	- retrieve a message or check if there is any message for the end user
	- check if a particular user is connected (e.g. who).
- Communications is via RPC. All errors are reported to the client.
- Server must implement:
	- connect()
	- disconnect()
	- deposit()
	- retrieve()
	- enquire() 
