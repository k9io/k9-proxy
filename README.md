
Join the Key9 Slack channel
---------------------------

[![Slack](./images/slack.png)](https://key9identity.slack.com/)


What is they Key9 Proxy?
------------------------

The Key9 proxy is a simple “caching” proxy for the Key9 API.  It is meant to provide high availability and accessibility to Key9 API users.  It acts as “middleware” to “cache” data between the Key9 API and Key9 users.  In an outage, the proxy can still service users via on-disk cache.

The proxy can be run as a network or host-based proxy.

Use cases:
----------

Some networks might have restrictions that prohibit machines from accessing the Key API directly.   In those cases,  the Key9 proxy can be used as a centralized API access point for all machines within a restricted network. 

As machines within your network make Key9 API requests,  the data is “caches”.   This cache services as a high availability mechanism in that, if the Key9 API is unreachable, data can be retrieved from the cache, thus allowing access to internal resources during the outage

The Key9 Proxy might be used as a local host proxy (cache) in certain situations. 

For example, if a host has unreliable Internet access, the Key9 proxy can keep a local cache of authentication data used by the operating system and SSH public keys.  This data can be refreshed when Internet access is reestablished.

What software uses they Key9 Proxy?
-----------------------------------

The proxy is used by k9-ssh (public key retrieval) and k9-nss (operating system NSS library)

Building and installing the Key9 Proxy
--------------------------------------

Make sure you have Golang installed! 

<pre>
$ go mod init k9-proxy
$ go mod tidy
$ go build
$ sudo mkdir -p /opt/k9/etc /opt/k9/bin
$ sudo cp etc/k9-proxy.yaml /opt/k9/etc
$ sudo cp k9-proxy /opt/k9/bin
$ sudo /opt/k9/bin/k9-proxy 	 # Run from the command line... Control C exits
$ sudo cp k9-proxy.service /etc/systemd/system
$ sudo systemctl enable k9-proxy
$ sudo systemctl start k9-proxy
</pre>

Prebuild Key9 proxy binaries
----------------------------

If you are unable to access a Golang compiler, you can download pre-built/pre-compiled binaries. These binaries are available for various architectures (i386, amd64, arm64, etc) and multiple operating systems (Linux, Solaris, NetBSD, etc).

You can find those binaries at: https://github.com/k9io/k9-binaries/tree/main/k9-proxy

You will need a copy of the 'k9-proxy' configuation file.  That is located at: 

https://github.com/k9io/k9-proxy/blob/main/etc/k9-proxy.yaml

