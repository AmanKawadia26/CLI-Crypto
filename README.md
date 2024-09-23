# CryptoTracker

CryptoTracker is a command-line application built in Go that allows users to track cryptocurrency prices, set alerts, and manage user profiles. The application features an intuitive user interface with enhanced visuals using the Fatih color library.



✨ Features

User Authentication: Login and Signup functionalities.

Admin and User Panels: Separate panels for admin and regular users.

View Top Cryptocurrencies: Easily access the top 10 cryptocurrencies.

Search Functionality: Look up specific cryptocurrencies by name.

Graphical Representation: 30-days graph will be displayed of the cryptocurrency

Price Alerts: Set custom price alerts for your favorite cryptocurrencies.

Unavailable Crypto Tracking: Request tracking for cryptocurrencies not currently supported.

Notification Functionality: When the user logs in then its notifictions pops up.

Add or Delete User: Admin can delete or add an user.

Delegation Functionality: Admin can escalate any user role.




📸 Snippets

Home Page: 

![image](https://github.com/user-attachments/assets/57dc05ca-8bf1-44ba-b46e-46dbd188bf0b)

User Page: 

![image](https://github.com/user-attachments/assets/20d8f935-a697-4818-968a-bacf82e6a61f)
![image](https://github.com/user-attachments/assets/5b112dd1-71c1-4875-bc23-d883ec84677e)
![image](https://github.com/user-attachments/assets/13428039-e2a0-4451-9b1d-739bab4cb9e1)

Admin Page: 

![image](https://github.com/user-attachments/assets/584a6b45-7527-48b2-885f-baaa6599430f)
![image](https://github.com/user-attachments/assets/9ad1896a-eda4-490e-b02c-ff19ed4bb852)







⚙️ Installation

Prerequisites
Ensure that you have Go installed on your system. You can download it from the official Go website.

Clone the Repository - 

git clone https://github.com/AmanKawadia26/CryptoTracker.git

cd CryptoTracker


Retrieve Dependencies - 

Before running the project, ensure you have all the necessary Go packages by running the following command:

go get -u


Make sure to include the following packages:

Fatih Color Library: This library is used to enhance the user interface with colored output in the terminal.

go get -u github.com/fatih/color





🚀 Running the Application

Once the dependencies are set up, you can run the application with:

go run main.go

This command will compile and execute the main.go file, starting the CryptoTracker application.
