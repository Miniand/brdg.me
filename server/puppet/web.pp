group {"puppet":
	ensure => present,
}

file { "/etc/init/boredga.me-web.conf":
	ensure => present,
	content => "env BOREDGAME_WEB_SERVER_ADDRESS=\":80\"
start on (local-filesystems and net-device-up IFACE=eth0)
stop on shutdown
exec /usr/bin/boredga.me-web 2>&1 >> /var/log/boredga.me-web",
}

file { "/etc/init.d/boredga.me-web":
	ensure => link,
	target => "/lib/init/upstart-job",
}

service { "boredga.me-web":
	ensure => running,
	provider => "upstart",
	require => [
		File["/etc/init/boredga.me-web.conf"],
		File["/etc/init.d/boredga.me-web"],
	],
}
