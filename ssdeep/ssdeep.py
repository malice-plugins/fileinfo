#!/usr/bin/env python
# -*- coding: utf-8 -*-

__author__ = 'blacktop'
__copyright__ = '''Copyright (C) 2013-2014 Josh "blacktop" Maine
                   This file is part of Malice - https://github.com/maliceio/malice
                   See the file 'LICENSE' for copying permission.'''

import ssdeep

ssdeep.hash_from_file('/etc/resolv.conf')

hash1 = ssdeep.hash('Also called fuzzy hashes, Ctph can match inputs that have homologies.')
hash2 = ssdeep.hash('Also called fuzzy hashes, CTPH can match inputs that have homologies.')

ssdeep.compare(hash1, hash2)

# pip install ssdeep
