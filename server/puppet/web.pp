group {"puppet":
	ensure => present,
}

file { "/etc/init/brdg.me-web.conf":
	ensure => present,
	content => "env BRDGME_WEB_SERVER_ADDRESS=\"brdg.me:80\"
start on (local-filesystems and net-device-up IFACE=eth0)
stop on shutdown
exec /usr/bin/brdg.me-web 2>&1 >> /var/log/brdg.me-web",
}

file { "/etc/init.d/brdg.me-web":
	ensure => link,
	target => "/lib/init/upstart-job",
}

service { "brdg.me-web":
	ensure => running,
	provider => "upstart",
	require => [
		File["/etc/init/brdg.me-web.conf"],
		File["/etc/init.d/brdg.me-web"],
	],
}
