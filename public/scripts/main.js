var isTweeting = false;

document.addEventListener("DOMContentLoaded", function (event) {
	document.getElementById("toggle").checked = false;
	getTweetStatus();

	document.getElementById("toggle").addEventListener("change", function() {
		console.log("Toggled");

		if (isTweeting) {
			stopTweeting();
		} else {
			startTweeting();
		}
	});
})

function startTweeting() {
	httpGetAsync("http://av-capture-linux.byu.edu:9000/tweeter/start", function() {
		isTweeting = true;
		console.log("Command to start tweeting sent");
	});
}

function stopTweeting() {
	httpGetAsync("http://av-capture-linux.byu.edu:9000/tweeter/stop", function() {
		isTweeting = false;
		console.log("Command to stop tweeting sent");
	});
}

function getTweetStatus() {
	httpGetAsync("http://av-capture-linux.byu.edu:9000/tweeter/status", function(status) {
		if (status == "false") {
			isTweeting = false;
		} else {
			isTweeting = true;
		}

		document.getElementById("toggle").checked = isTweeting;
	});
}

function httpGetAsync(theUrl, callback) {
	var xmlHttp = new XMLHttpRequest();
	xmlHttp.onreadystatechange = function() {
		if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
			callback(xmlHttp.responseText);
	}

	xmlHttp.open("GET", theUrl, true); // true for asynchronous
	xmlHttp.send(null);
}
