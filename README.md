Someone was asking what to do about a thing not recognizing \r as a newline
sequence, so I made a thing that changes \r to \n unless the EOL scheme is
\r\n. Just wrap your reader in nl.New() and enjoy.