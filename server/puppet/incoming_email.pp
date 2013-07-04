group {"puppet":
	ensure => present,
}

file { "/etc/init/boredga.me-email.conf":
	ensure => present,
	content => "env BOREDGAME_EMAIL_SERVER_ADDRESS=\":81\"
start on (local-filesystems and net-device-up IFACE=eth0)
stop on shutdown
exec /usr/bin/boredga.me-email 2>&1 >> /var/log/boredga.me-email",
}

file { "/etc/init.d/boredga.me-email":
	ensure => link,
	target => "/lib/init/upstart-job",
}

service { "boredga.me-email":
	ensure => running,
	provider => "upstart",
	require => [
		File["/etc/init/boredga.me-email.conf"],
		File["/etc/init.d/boredga.me-email"],
	],
}
