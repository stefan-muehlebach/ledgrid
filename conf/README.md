# The layout of your LEDGrid

## Modules

The hardware I use for my projects is based on chains of WS2801 NeoPixels.
You can buy them in strands of 25 or 50 pixels. A chain with 100 pixels,
organized in a square of 10 by 10 pixels is the smallest building block, I
call it a _module_. In the picture below, you'll see how the 100 pixels are
organized in such a module.

![LR:000](plots/LR-000.png)

This particular module is called **LR:000** (more on this name later).

Each circle stands for one pixel (i.e. LED). The numbers in the circles as
well as the big gray strand in the background show the direction of the chain.
The green circle denotes the first pixel of the module whereas the red circle
denotes the last pixel of the module. Because the chain starts on the left
side of the module and ends on the right side, this type of module is called
**LR** (left-to-right).

The mirrored twin of a **LR** module is (of course) a **RL** module (shown in
the picture below):

![RL:000](plots/RL-000.png)

Each of the two modules can be turned by 0, 90, 180 or 270 degrees, which
gives a total of 8 modules

| LR:000 | LR:090 | LR:180 | LR:270 |
|--------|--------|--------|--------|
|![LR:000](plots/LR-000.png)|![LR:090](plots/LR-090.png)|![LR:180](plots/LR-180.png)|![LR:270](plots/LR-270.png)|

| RL:000 | RL:090 | RL:180 | RL:270 |
|--------|--------|--------|--------|
|![RL:000](plots/RL-000.png)|![RL:090](plots/RL-090.png)|![RL:180](plots/RL-180.png)|![RL:270](plots/RL-270.png)|

## Chaining the modules

With this set of 8 modules you can start and build your LED-wall. There is no
limit in size and you dont't have to place the modules seamingless - as long
as you respect the following _golden rule_:

> The red end circle of one module must be in the same corner as the green
> start circle of the following module.

400 pixels, mounted on 4 modules can be arranged in the following orders.
For each example, the content of it's configuration file is also shown.
The syntax of the configuration files will be explained right afterwards.

### Adjacent

The simplest configuration: just one module after the other.

![](plots/sample01.png)

    [
        {"Col": 0, "Row": 0, "Mod": "LR:0"  },
        {"Col": 1, "Row": 0, "Mod": "LR:0"  },
        {"Col": 2, "Row": 0, "Mod": "LR:0"  },
        {"Col": 3, "Row": 0, "Mod": "LR:0"  }
    ]

### Even-Odd

A configuration with a lot of space between the modules.
Note that the _golden rule_ is still valid: red and green circles meet in one
corner.

![](plots/sample02.png)

    [
        {"Col": 0, "Row": 1, "Mod": "LR:0"  },
        {"Col": 1, "Row": 0, "Mod": "RL:180"},
        {"Col": 2, "Row": 1, "Mod": "LR:0"  },
        {"Col": 3, "Row": 0, "Mod": "RL:180"}
    ]

### Odd-Even

Same as above - but different ;-)

![](plots/sample03.png)

    [
        {"Col": 0, "Row": 0, "Mod": "RL:180"},
        {"Col": 1, "Row": 1, "Mod": "LR:0"  },
        {"Col": 2, "Row": 0, "Mod": "RL:180"},
        {"Col": 3, "Row": 1, "Mod": "LR:0"  }
    ]

### A Hole

You can event arrange the modules in a circle around a square hole.

![](plots/sample04.png)

    [
        {"Col": 1, "Row": 0, "Mod": "RL:180"},
        {"Col": 2, "Row": 1, "Mod": "RL:90" },
        {"Col": 1, "Row": 2, "Mod": "RL:0"  },
        {"Col": 0, "Row": 1, "Mod": "RL:270"}
    ]

## Big walls and the Default configuration

In most cases, you just want a seamingless configuration, without holes and
irregularities. In this case, a default configuration can be choosen, which
automaticaly places the modules in order to fill the specified area.

The final size of the grid must (of course) be a multiple of 10.
The grids below are all automatically generated.

### 40 x 10 pixels

400 pixels, the length of the wire from first to last led is 40m

![](plots/default40x10.png)

### 40 x 40 pixels

1600 pixels, the length of the wire from first to last led is 160m

![](plots/default40x40.png)

### 80 x 60 pixels

4800 pixels, the length of the wire from first to last led is 480m - and
probably just a dream.

![](plots/default80x60.png)
