# DataTable

Play with csv's on a terminal

# Pre-requisites

- Go 1.12+ 

# Setup
    
    make datatable

# Test
    
    make test
    
# Run UI

    ./datatable ui -f <csv filename>
    
It will read the file passed and show it in the ui.
Use your arrow keys to choose different entries in the csv.
Once you have chosen the entry you are interested in, type `e` to change it.

The item should appear in the bottom window.

You can change it to the value you want by typing and `backspace` and when you are ready press `enter`
to edit the value.

You can use the arrow keys while editting the value.

Finally, after making all the changes you want, type `s` to save the changes back to the csv.

Type `w` to save as.

press `ctrl-h` at any time  to show the above help.

press `ctrl-c` to exit
