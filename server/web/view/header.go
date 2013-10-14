package view

var headerTmpl = `<html>
<head>
<link rel="stylesheet" type="text/css" href="//cdnjs.cloudflare.com/ajax/libs/pure/0.3.0/pure-min.css">
<style type="text/css">
html,
button,
input,
select,
textarea,
.pure-g,
.pure-g [class *= "pure-u"],
.pure-g-r,
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
	float: left;
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
	float: right;
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
</style>
<title>brdg.me - {{.}}</title>
<body>
<header>
	<div class="logo"><a href="/">{{template "title"}}</a></div>
	<div class="account">
		<a href="/sign-in">Sign in / register</a>
	</div>
	<div class="header-end"></div>
</header>`
