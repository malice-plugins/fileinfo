#!/usr/bin/env python
# -*- coding: utf-8 -*-

__author__ = 'blacktop'
__copyright__ = '''Copyright (C) 2013-2016 Josh "blacktop" Maine
                   This file is part of Malice - https://github.com/maliceio/malice
                   See the file 'LICENSE' for copying permission.'''

import envoy


class TrID():
    name = "TrID"
    description = "TrID is an utility designed to identify file types " \
                  "from their binary signatures. While there are similar " \
                  "utilities with hard coded logic, TrID has no fixed rules. " \
                  "Instead, it's extensible and can be trained to recognize " \
                  "new formats in a fast and automatic way."
    severity = 0
    categories = ["file type"]
    authors = ["blacktop"]
    references = ["http://mark0.net/soft-trid-e.html"]
    minimum = "v0.1.0-alpha"

    def __init__(self, path):
        self.path = path

    @staticmethod
    def format_output(output):
        trid_results = []
        results = output.split('\n')
        results = filter(None, results)
        for trid in results:
            trid_results.append(trid)
        return trid_results

    @staticmethod
    def update_definitions():
        # Update the TRiD definitions
        try:
            r = envoy.run('python /opt/info/trid/tridupdate.py', timeout=20)
            return r.std_out
        except AttributeError:
            print 'ERROR: TrID Failed.'
            return 'trid', dict(error='TrID failed to run.')
        except Exception, e:
            print e.message

        return None

    def scan(self):
        # Run exiftool on tmp file
        try:
            r = envoy.run('/opt/info/trid/trid ' + self.path, timeout=15)
        except AttributeError:
            print 'ERROR: TrID Failed.'
            return 'trid', dict(error='TrID failed to run.')
        except Exception, e:
            return dict(exif=dict(error=e.message))

        # return key, stdout as a dictionary
        return dict(trid=self.format_output(r.std_out.split(self.path)[-1]))
