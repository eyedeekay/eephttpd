eepHttpd - easy-to-use tool for setting up I2P Sites
====================================================

eepHttpd is the easiest tool for hosting new sites on the I2P network from your own
computer(As long as they are mostly *static* sites, for now). You start the application
and it sets up a directory in a logical place on your system.

Homepages:
----------

- **[I2P Site](http://tvndxxkxcstbtqfxg7iigco6bj22ff2y6jxikmk7wqkyadkhrd4a.b32.i2p/)**
- **[Visible Internet Site](https://eyedeekay.github.io/eephttpd)**

Screenshot:
-----------

![Initial Setup Screen](eephttpd.png)

Example Setup:
--------------

On Windows the default site will be created in the "My Documents" folder, and it will be
empty by default. For example:

        C:\\Documents and Settings\User\My Documents\www

On Unixes(Linux and OSX are tested but any Unix should work) it will be the "www" directory
in the directory where you run the application. So if you ran the application from

        /home/user/eephttpd/

then you would end up with a `www` directory there, for example:

        /home/user/eephttpd/www/

Just put the files you want to serve, like your web site or open directory of content, inside
of that directory.

Example Clone:
--------------

eepHttpd is capable of mirroring a static site stored in a git repository. Sort of like
a self-hosted github page. This feature is accessible via the GUI. To do this, fill in the
`Clone Site from a git repository` section.

![Initial Clone Screen](eephttpd-clone.png)

This example will mirror the eephttpd site itself. It should work for most github pages and
all static sites.

Learn More:
-----------

- **[Source Code](https://github.com/eyedeekay/eephttpd)**
- **[File an issue, request a feature](https://github.com/eyedeekay/eephttpd/issues)**
- **[File an issue, request a feature](https://github.com/eyedeekay/eephttpd/pulls)**
