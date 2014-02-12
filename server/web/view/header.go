package view

var headerTmpl = `<html>
<head>
<link rel="stylesheet" type="text/css" href="//cdnjs.cloudflare.com/ajax/libs/pure/0.3.0/pure-min.css">
<style type="text/css">
body,
.pure-g [class *= "pure-u"],
.pure-g-r [class *= "pure-u"] {
    font-family: monospace;
}

a:visited {
	color: blue;
}

.title {
	margin-top: 4em;
	text-align: center;
}

.title h1 {
	text-align: left;
	display: inline-block;
	white-space: pre-wrap;
}

header {
	border-bottom: 1px solid black;
	padding: 1em;
}

header .logo {
	position: absolute;
	left: 2em;
	top: 2em;
	height: 4em;
	text-align: left;
	display: inline-block;
	white-space: pre-wrap;
	font-size: 0.5em;
	font-weight: bold;
}

header .logo a {
	text-decoration: none;
}

header .account {
	position: absolute;
	right: 2em;
	top: 1em;
	height: 4em;
	line-height: 4em;
}

header .header-end {
	clear: both;
}

.game-list {
	max-width: 960px;
	margin: 5em auto;
}

.game-list .game {
	text-align: center;
}

.active-games {
	margin-left: 14em;
	margin-right: 14em;
	margin-top: 0.5em;
	white-space: nowrap;
	overflow-x: auto;
}

.active-game {
	display: inline-block;
	white-space: normal;
	line-height: 1.5em;
	border: 1px dashed gray;
	padding: 0 1em;
	margin-bottom: 0.5em;
}

.active-game-your-turn {
	font-weight: bold;
	border: 1px solid gray;
	background-color: silver;
}

.game-show {
	margin: 0 auto;
	max-width: 100em;
}

.active-game-name {
	text-align: center;
}

.active-game-player {
	text-align: center;
}

.game-output-container {
	text-align: center;
}

.game-output {
	text-align: left;
	display: inline-block;
	white-space: pre-wrap;
	font-size: 1.3em;
}

.game-log-container {
	text-align: center;
}

.game-log {
	text-align: left;
	white-space: pre-wrap;
	max-height: 50em;
	overflow-y: auto;
}

.game-log-heading {
	text-align: center;
}

.game-log {
	white-space: pre-wrap;	
}

.game-input-container {
	text-align: center;
}

.game-input-available-commands {
	display: inline-block;
	text-align: left;
	margin-bottom: 1em;
	white-space: pre-wrap;	
}

.game-input-command {
	width: 70%;
}

.game-input-submit {
	width: 20%;
}
</style>
<title>brdg.me - {{.}}</title>
</head>
<script src="https://login.persona.org/include.js"></script>
<script>
function getCookie(name) {
	var value = "; " + document.cookie;
	var parts = value.split("; " + name + "=");
	if (parts.length == 2) return parts.pop().split(";").shift();
}
var loggedInUser = "{{loggedInUser}}";
if (loggedInUser === "") {
	loggedInUser = null;
}
navigator.id.watch({
	loggedInUser: loggedInUser,
	onlogin: function(assertion) {
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/session/sign-in", true);

		xhr.onreadystatechange = function() {
			if (xhr.readyState == 4) {
				if (xhr.status == 200) {
					// reload page to reflect new login state
					window.location.reload();
				} else {
					navigator.id.logout();
				}
			}
		};

		var fd = new FormData();
		fd.append("assertion", assertion)

		xhr.send(fd);
	},
	onlogout: function() {
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/session/sign-out", true);

		xhr.onreadystatechange = function() {
			if (xhr.readyState == 4) {
				if (xhr.status == 200) {
					// reload page to reflect new login state
					window.location.reload();
				} else {
					navigator.id.logout();
				}
			}
		};

		xhr.send();
	}
});
var session = getCookie('session');
if (session) {
	var ws = new WebSocket("ws://" + window.location.host + "/ws");
	ws.onopen = function() {
		ws.send(session);
	};
	ws.onmessage = function(event) {
		console.log(event.data);
	};
	ws.onclose = function() {
		alert("Your connection to the server has dropped, please refresh the page");
	};
	window.ws = ws;
}
</script>
<body>
<header>
	<div class="logo"><a href="/">{{template "title"}}</a></div>
	<div class="active-games">
		<a href="/" class="active-game active-game-your-turn">
			<div class="active-game-name">Texas hold 'em (183f94ca)</div>
			<div class="active-game-player">Your turn</div>
		</a> <a href="/" class="active-game active-game-your-turn">
			<div class="active-game-name">Texas hold 'em (183f94ca)</div>
			<div class="active-game-player">Your turn</div>
		</a> <a href="/" class="active-game">
			<div class="active-game-name">Texas hold 'em (183f94ca)</div>
			<div class="active-game-player">baconheist's turn</div>
		</a> <a href="/" class="active-game">
			<div class="active-game-name">Texas hold 'em (183f94ca)</div>
			<div class="active-game-player">striker203 and ashtermet's turn</div>
		</a> <a href="/" class="active-game">
			<div class="active-game-name">Texas hold 'em (183f94ca)</div>
			<div class="active-game-player">baconheist's turn</div>
		</a> <a href="/" class="active-game">
			<div class="active-game-name">Texas hold 'em (183f94ca)</div>
			<div class="active-game-player">baconheist's turn</div>
		</a> <a href="/" class="active-game">
			<div class="active-game-name">Texas hold 'em (183f94ca)</div>
			<div class="active-game-player">baconheist's turn</div>
		</a>
	</div>
	<div class="account">
		{{if loggedInUser}}
		<a href="javascript:navigator.id.logout()">Sign out</a>
		{{else}}
		<a href="javascript:navigator.id.request()">Sign in</a>
		{{end}}
	</div>
	<div class="header-end"></div>
</header>`
