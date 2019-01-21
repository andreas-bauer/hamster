[![Build Status](https://travis-ci.org/michaeldorner/Gerry.svg)](https://travis-ci.org/michaeldorner/Gerry) 
[![codecov](https://codecov.io/gh/michaeldorner/Gerry/branch/master/graph/badge.svg)](https://codecov.io/gh/michaeldorner/Gerry)
[![codebeat badge](https://codebeat.co/badges/f8306b22-3837-4244-a637-e880c6532700)](https://codebeat.co/projects/github-com-michaeldorner-gerry-master)
![license](https://img.shields.io/github/license/mashape/apistatus.svg)

# Hamster

> Hamster is a framework for crawl and store REST API responses written in Go. 


## Table of Contents

- [Installation](#install)
- [Usage](#usage)
- [Contributions](#contributions)
- [License](#license)


## Installation



## Usage

    python gerry.py <gerrit_instance> [--directory=<storage_directory>]
    
* `<gerrit_instances>`: Gerry supports the gerrit instances of OpenStack (`openstack`), Chromium (`chromium`), Gerrit (`gerrit`), Android ('android'), Go (`golang`), LibreOffice (`libreoffice`), Eclipse (`eclipse`), Wikimedia (`wikimedia`), and ONAP (`onap`). 
* `<storage_directory>` (optional): The storage directory where to store the files (default `./gerry_data/`).


## Contributions

### To-do

It would be great to get a pull request containing new instances or adding the option to add own (private, non-open-source) instances. 

Other pull requests are always welcome. 


### Acknowledgements

Many thanks to my excellent master student [Jonathan Frieß](https://github.com/FreezerJohn), who added the full test setup. 


## License 

Gerry is released under the MIT license. See [LICENSE](LICENSE) for more details.
