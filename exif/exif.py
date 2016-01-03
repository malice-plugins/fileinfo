#!/usr/bin/env python
# -*- coding: utf-8 -*-

__author__ = 'Josh Maine'
__copyright__ = '''Copyright (C) 2013-2014 Josh "blacktop" Maine
                   This file is part of Malice - https://github.com/blacktop/malice
                   See the file 'docs/LICENSE' for copying permission.'''

import tempfile
from os import unlink

import envoy
from lib.common.abstracts import FileAnalysis
from lib.common.out import print_error

# from os.path import exists

ignore_tags = ['Directory',
               'File Name',
               'File Permissions',
               'File Modification Date/Time']


class Exif(FileAnalysis):

    name = "ExifTool"
    description = "ExifTool is a platform-independent Perl library plus a command-line " \
                  "application for reading, writing and editing meta information in a " \
                  "wide variety of files."
    severity = 2
    categories = ["file type"]
    authors = ["blacktop"]
    references = ["http://www.sno.phy.queensu.ca/~phil/exiftool/"]
    minimum = "v0.1-alpha"
    # evented = True

    def __init__(self, data):
        FileAnalysis.__init__(self, data)
        self.data = data

    def format_output(self, output):
        exif_tag = {}
        exif_results = output.split('\n')
        exif_results = filter(None, exif_results)
        for tag in exif_results:
            tag_part = tag.split(':', 1)
            if len(tag_part) == 2:
                if tag_part[0].strip() not in ignore_tags:
                    exif_tag[tag_part[0].strip()] = tag_part[1].strip().decode('utf-8')
        return exif_tag

    def scan(self):
        # : create tmp file
        handle, name = tempfile.mkstemp(suffix=".data", prefix="exif_")
        #: Write data stream to tmp file
        with open(name, "wb") as f:
            f.write(self.data)
        #: Run exiftool on tmp file
        try:
            r = envoy.run('exiftool ' + name, timeout=15)
        except AttributeError:
            print_error('ERROR: Exif Failed.')
            return 'exif', dict(error='Exiftool failed to run.')
        else:
            #: return key, stdout as a dictionary
            return 'exif', self.format_output(r.std_out)
        finally:
            #: delete tmp file
            unlink(name)
            # exists(name)
