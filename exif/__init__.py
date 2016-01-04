#!/usr/bin/env python
# -*- coding: utf-8 -*-

__author__ = 'blacktop'
__copyright__ = '''Copyright (C) 2016 Josh "blacktop" Maine
                   This file is part of Malice - https://github.com/maliceio/malice
                   See the file 'LICENSE' for copying permission.'''

import envoy

ignore_tags = ['Directory', 'File Name', 'File Permissions', 'File Modification Date/Time']


class Exif():
    name = "ExifTool"
    description = "ExifTool is a platform-independent Perl library plus a command-line " \
                  "application for reading, writing and editing meta information in a " \
                  "wide variety of files."
    severity = 2
    categories = ["file type"]
    authors = ["blacktop"]
    references = ["http://www.sno.phy.queensu.ca/~phil/exiftool/"]
    minimum = "v0.1.0-alpha"

    def __init__(self, path):
        self.path = path

    @staticmethod
    def format_output(output):
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
        #: Run exiftool on file
        try:
            r = envoy.run('exiftool ' + self.path, timeout=15)
        except AttributeError:
            print 'ERROR: Exif Failed.'
            return 'exif', dict(error='Exiftool failed to run.')
        except Exception, e:
            print e.message
        else:
            #: return key, stdout as a dictionary
            return dict(exif=self.format_output(r.std_out))
