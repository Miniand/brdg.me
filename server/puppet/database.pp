group {"puppet":
	ensure => present,
}

package {"mongodb":
	ensure => present,
}
