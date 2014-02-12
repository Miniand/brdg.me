group {"puppet":
	ensure => present,
}

file { "/etc/init/brdg.me.conf":
	ensure => present,
	content => "env BRDGME_WEB_SERVER_ADDRESS=\"brdg.me:80\"
env BRDGME_EMAIL_SERVER_ADDRESS=\":81\"
start on (local-filesystems and net-device-up IFACE=eth0)
stop on shutdown
exec /usr/bin/brdg.me 2>&1 >> /var/log/brdg.me",
}

file { "/etc/init.d/brdg.me":
	ensure => link,
	target => "/lib/init/upstart-job",
}

service { "brdg.me":
	ensure => running,
	provider => "upstart",
	require => [
		File["/etc/init/brdg.me.conf"],
		File["/etc/init.d/brdg.me"],
	],
}
