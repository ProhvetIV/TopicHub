import { EventBus } from '@/eventBus';
import { getWebSocketService } from '@/websocket.js';

// input checks
function isStringEmptyOrWhitespace(str) {
    return !str.replace(/\s+/g, ''); // Returns true if string is empty or contains only whitespace characters
}

function isVariableEmpty(variable) {
    return variable === undefined || variable === null;
}

export function checkIfFilled(obj) {
    for (let key in obj) {
        let value = obj[key];
        if (typeof value === 'string') {
            if (isStringEmptyOrWhitespace(value)) {
                alert(`String '${key}' is empty or contains only whitespace characters.`);
                console.log(`String '${key}' is empty or contains only whitespace.`);
                return true
            }
        } else {
            if (isVariableEmpty(value)) {
                alert(`Variable '${key}' is empty.`);
                console.log(`Variable '${key}' is empty.`);
                return true
            }
        }
    }
    return false
}

// Cookies!

export function createSession(wsMessageData) {
	const username = wsMessageData.username;
	const sessionID = wsMessageData.sessionID;
	const expirationTime = new Date(wsMessageData.expirationTime);
	sessionStorage.setItem("username", username);
	sessionStorage.setItem("sessionID", sessionID);
	sessionStorage.setItem("expirationTime", expirationTime.getTime().toString());
    EventBus.usrn = username;
}

export function isSessionExpired() {
	const expirationTime = parseInt(sessionStorage.getItem("expirationTime"));
	if (expirationTime) {
		const now = Date.now();
		return now > expirationTime;
	}
	return true; // Session expiration time not found
}

export function removeSession() {
	sessionStorage.removeItem("sessionID");
	sessionStorage.removeItem("expirationTime");
	sessionStorage.removeItem("username");
    EventBus.usrn = null;
}

export function getTimeSince(date) {
	var creationDate = new Date(date);
	var currentDate = new Date();

	var timeDifference = currentDate - creationDate;
	var secondsDifference = Math.floor(timeDifference / 1000);

	// Convert seconds to minutes, hours, and days
	var minutes = Math.floor(secondsDifference / 60); //
	var hours = Math.floor(minutes / 60);
	var days = Math.floor(hours / 24);

	if (days != 0) {
		return days + "d " + (hours % 24) + "h";
	} else if (hours != 0) {
		return (hours % 24) + "h " + (minutes % 60) + "min";
	} else {
		return (minutes % 60) + "min";
	}
}

export function formatDate(date){
	let newDate = new Date(date);
  	let formatDate = newDate.toLocaleString();
	return formatDate;
}

export function checkUser(sender, reciever) {
	let recieverUsername = document.querySelector(".chatUsername").textContent;
	if (sender === recieverUsername || recieverUsername === reciever) {
		return true;
	}
	return false;
}

export const sendMessageWhenOpen = (type, data) => {
	const webSocketService = getWebSocketService();
	const interval = setInterval(() => {
        if (webSocketService.socket.readyState === WebSocket.OPEN) {
            clearInterval(interval);
            webSocketService.sendMessage(type, data);
        }
    }, 100);
};
let count = 0

export function showMessages(message, before = true) {
	let username = document.querySelector(".chatUsername");
	const headerUsername = sessionStorage.getItem('username')

	if (
		(count > 0 && message.senderUsername === headerUsername) ||
		message.senderUsername === username.textContent ||
		message.recieverUsername === username.textContent
	) {
		const messagesDiv = document.querySelector(".messages");
		const messageData = document.createTextNode(message.content);
		const senderName = document.createTextNode(message.senderUsername + ": ");
		const messageDateTime = document.createTextNode(formatDate(message.creation_date));
		// adds space between each text item
		let spanElement = document.createElement("span");
		spanElement.classList.add("highlighted")
		spanElement.appendChild(messageData);

		let spanDateTime = document.createElement("span");
		spanDateTime.classList.add("dateTime")
		spanDateTime.appendChild(messageDateTime);

		let br = document.createElement("br");
		let singleMessage = document.createElement("div");
		singleMessage.classList.add("singleMessage")
		singleMessage.appendChild(senderName);
		singleMessage.appendChild(spanElement);
		singleMessage.appendChild(spanDateTime);
		singleMessage.appendChild(br);
		if (before) {
			messagesDiv.insertBefore(singleMessage, messagesDiv.firstChild);
		} else {
			messagesDiv.appendChild(singleMessage, messagesDiv.firstChild);
		}

		scrollToBottom(".messages");
	}
}

function scrollToBottom(selector) {
	const messagesContainer = document.querySelector(selector);
	messagesContainer.scrollTop = messagesContainer.scrollHeight;
}

export function notification(data) {
	const users = document.querySelectorAll(".singleUserButton button");
	const usersContainer = document.querySelector(".Users");

	users.forEach((user) => {
		let singleUserButton = user.closest(".singleUserButton"); // Find the closest parent div with class 'singleUserButton'
		let spanCounter = singleUserButton.querySelector(".notification");
		if (user.textContent === data.senderUsername) {
			let counter = parseInt(spanCounter.textContent || 0);
			counter++;
			spanCounter.textContent = counter;
			spanCounter.classList.add("circle");
			usersContainer.insertBefore(singleUserButton, usersContainer.firstChild);
		}
	});
}

export const getMimeType = (base64) => {
	if (base64.startsWith('/9j/')) return "image/jpeg" // Works for both jpeg and jpg
	if (base64.startsWith('iVBORw0KGgo')) return "image/png"
	if (base64.startsWith('R0lGODdh') || base64.startsWith('R0lGODlh')) return 'image/gif';
	if (base64.startsWith('UklGR')) return 'image/webp';
	return 'application/octet-stream';
}

export const emojis = ["😄", "😍", "😎", "🤗", "😇", "😊", "😋", "😘", "🥰", "😜", 
"🤩", "🥳", "😷", "🤠", "🥺", "😢", "😡", "🤬", "😱", "😴", "💀", '😈', '👻', '🤖', '👽', '🎃', '💩'];

