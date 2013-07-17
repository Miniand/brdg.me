group {"puppet":
	ensure => present,
}

file { "/etc/init/brdg.me-email.conf":
	ensure => present,
	content => "env BRDGME_EMAIL_SERVER_ADDRESS=\":81\"
start on (local-filesystems and net-device-up IFACE=eth0)
stop on shutdown
exec /usr/bin/brdg.me-email >> /var/log/brdg.me-email 2>&1",
}

file { "/etc/init.d/brdg.me-email":
	ensure => link,
	target => "/lib/init/upstart-job",
}

service { "brdg.me-email":
	ensure => running,
	provider => "upstart",
	require => [
		File["/etc/init/brdg.me-email.conf"],
		File["/etc/init.d/brdg.me-email"],
	],
}
