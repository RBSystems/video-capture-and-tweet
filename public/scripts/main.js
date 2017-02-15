var isTweeting = false;

document.addEventListener("DOMContentLoaded", function (event) {
	document.querySelector("input[name=toggle]").addEventListener("change", function() {
		console.log("Toggled");

		if (isTweeting) {
			stopTweeting();
		} else {
			startTweeting();
		}
	});
})

function startTweeting() {
	httpGetAsync("http://av-capture-linux.byu.edu/tweeter/start");
}

function stopTweeting() {
	httpGetAsync("http://av-capture-linux.byu.edu/tweeter/stop");
}

function getTweetStatus() {
	httpGetAsync("http://av-capture-linux.byu.edu/tweeter/status", function(status) {
		isTweeting = status;
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
