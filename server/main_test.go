package main

import (
	"testing"
)

func TestReceiveMail(t *testing.T) {
	_ = "Received: from beefsack.com (localhost [127.0.0.1])	by boredga.me (Postfix) with SMTP id 6D9EC40625	for <play@boardga.me>; Thu,  4 Jul 2013 07:23:28 +0000 (UTC)Subject: Another test to play userMessage-Id: <20130704072342.6D9EC40625@boredga.me>Date: Thu,  4 Jul 2013 07:23:28 +0000 (UTC)From: beefsack@beefsack.comPlay user!  Here's another test for you."
}
