Prometheus exporter for ``libvirtd``
====================================
Prometheus exporter for Libvirt metrics, currently exposing domain stats
at the moment but with the ability for pluggable metric collectors.


Building
--------
You can build the project yourself locally simply by running the following
inside the root folder.  You'll need to make sure that you have the headers
for ``libvirtd`` installed on your machine.  The following example is for
a Debian based machine.

.. code-block:: bash

   apt-get -y install libvirt-dev
   go build


Usage
-----
There are a few different ways that you can choose to deploy this exporter,
it's up to you to choose which one you prefer.

Docker
~~~~~~
``vexxhost/libvirtd_exporter:latest`` always points at the latest tested
commit which is always gated so it should not break and you can rely on
deploying it.  When running with Docker, you'll need to mount the ``libvirt``
socket into the container, preferebly the read-only one.


Contributing
------------

Running Locally
~~~~~~~~~~~~~~~
There are scenarios where you need to iterate on the code lcoally but run it
against a remote hypervisor.  It's possible to do this over SSH, an example
of how to do this against a CentOS host with ``libvirtd`` is:

.. code-block:: bash

   go run libvirtd_exporter.go --libvirt.uri="qemu+ssh://root@remote-system/system?socket=/var/run/libvirt/libvirt-sock-ro"
