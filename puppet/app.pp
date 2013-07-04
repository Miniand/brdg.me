group {"puppet":
	ensure => present,
}

file { "/etc/init/boredga.me.conf":
	ensure => present,
	content => "env SERVER_ADDRESS=\"localhost:80\"
start on (local-filesystems and net-device-up IFACE=eth0)
stop on shutdown
exec /usr/bin/boredga.me 2>&1 >> /var/log/boredga.me",
}

file { "/etc/init.d/boredga.me":
	ensure => link,
	target => "/lib/init/upstart-job",
}

service { "boredga.me":
	ensure => running,
	provider => "upstart",
	require => [
		File["/etc/init/boredga.me.conf"],
		File["/etc/init.d/boredga.me"],
	],
}
