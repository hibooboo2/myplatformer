Tile Texture Naming Format
Updated 5/1/97

[AA][BB][C][D][E][F]
------------------

  Group Type 1

[BB]
Group Type Transition To
If center tile - xy location if this is used.  (Unused characters for
artist.)
[C]
  Tile Type
         a = 00
             11
 
         b = 10
 
         c = 11
             00
 
         d = 01
             01
         e = 00
             10
 
         f = 10
             00
 
             00
 
         h = 00
             01
 
         i = 10
 
         j = 11
             10
 
         k = 11
             01
 
         l = 01
             11

         m = 00
             00
The main tile. Tile number refers to varients on the main tile that may be swapped in without any noticable side effects. These varients can be used to tile large areas. 

         n = Single tile details.  Any textures that tile with the main tiles but are not meant to be used over large areas. 

         s = Special large details.  Any textures that are larger than a single tile.  Typically part of a main tile set.  The version tile numbers are the tile’s position in the brush in x,y coordinates. Note that this allows us to have brushes of dimensions upto 36 by 36 which should be more than enough. The BB numbers are used to label different brushes. Most typically manipulated as brushes in the terrain editor. 

[D]
  Color Depth
  8 - 256 colors
  T - Truecolor

[E]     
  Tile version type.  For example, two main tiles might have two different borders connecting them.  The different borders would also have different version numbers. 

[F]
  Random modification on tile that Auto randomizer uses when selecting tiles.   

*Main Tile List*

B0 = Blob 0
B1 = Blob 1
D0 = Dirt 0 
D1 = Dirt 1
G0 = Grass 0
R0 = Ramp 0 
S0 = Stone 0
S1 = Stone 1
S2 = Stone 2
S3 = Stone 3 (Squiggle Rock)
S4 = Stone 4 (Hex pavement)

