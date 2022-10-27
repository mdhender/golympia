# Olympia Map generator 1.0
## Instructions for Use

### Start with an ASCII map
Construct an ascii character map of the world.
Each character represents one province.
Use the following characters to vary terrain:

    . , \ space    ocean
    m M            mountain
    d D            desert
    p P            plain
    f F            forest
    ; ~ " :        sea lane
    s S            swamp
    -              steppe
    o              random (30% forest/plain, 20% mountain, 10% swamp/desert)
    ?              hidden province
    #              impassable region (will be a "hole" in the map)
    * %            city

### Hidden Provinces
Take care when placing `?` provinces.
The terrain type will be inferred from one of the neighboring squares.

### Create a Regions file
Create a Regions file of the form:

    AA    Name of region AA
    AB    Name of region AB
    aa    Name of ocean aa
    zz    Name of ocean zz

Do not use the letter `o` in ocean/sea codes, since `o` represents plains on the map.
I would have used another character besides `o` for plains, but none looked as good visually.

### Name Regions and Oceans
Name regions and oceans by placing the two-letter code for them somewhere in the region.
Again, take care when placing the two-letter codes,
as the terrain type of the two squares taken up by the code will be inferred from neighboring provinces.

### Delineate Seas
Delineate seas by varying the character used to represent them.
In this way, the first sea name will not flood-fill the entire ocean, since it will be stopped at the border.

For example:

        +-------------------------+
        |~~~~~~........    ,,,,,,,|
        |~~~~~~~~...sa..    ,,,,,,|
        |~~~sb~~~~...... as ,,cs,,|
        |~~~~~~~~~......     ,,,,,|
        +-------------------------+

        sa    Sea of Athens
        sb    Sea of Boetia
        as    Aedras Sea
        cs    Carthas Sea

Spaces, `~`, `,`, and `.` may be used to represent seas.

### Generate the map
Run the map generator:

    goly generate map --map-file _map_name_ --loc-file _loc_file_name_

The _loc_file_name_ is the Olympia database `libdir/loc` file.
The _map_name_ should be your ASCII character map file.
The output will contain useful information about land continents as well as any warnings issued during map generation.

### Map Wrap
Although the Olympia map will wrap at the edges, for the purposes of flood-fill algorithms the map generator will stop at the edges.

### Named Cities
A `*` may be placed on land regions to indicate that a named city should go there.
Additional random cities will be scattered across the map.
City names will be read from the "Cities" file.
