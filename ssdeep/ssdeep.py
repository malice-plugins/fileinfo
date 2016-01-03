import ssdeep

ssdeep.hash_from_file('/etc/resolv.conf')

hash1 = ssdeep.hash('Also called fuzzy hashes, Ctph can match inputs that have homologies.')
hash2 = ssdeep.hash('Also called fuzzy hashes, CTPH can match inputs that have homologies.')

ssdeep.compare(hash1, hash2)

# pip install ssdeep
