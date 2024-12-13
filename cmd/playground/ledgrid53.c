//go:build ignore
// +build ignore

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <signal.h>
#include <limits.h>
#include <string.h>
#include <ncurses.h>
#include <assert.h>
#include <math.h>
#include <libgen.h>

#include "PiPack.h"

//-----------------------------------------------------------------------------
//
// ledgrid53 --
//
//
//     Plasma Effekt Generator.
//
//-----------------------------------------------------------------------------

// PaletteType --
//
//     Definition einer Farbpalette mit den Stuetzfarben (colorList) und den
//     berechneten Zwischenfarben (shadeList).
//
typedef struct PaletteType {
    char *name;
    int numColors;
    unsigned int *colorList;
    unsigned char *shadeList;
} *PaletteType;

PaletteType Pal_Init (FILE *fp) {
    PaletteType pal;
    char palName[256];
    unsigned int colorList[128];
    int numColors, i;

    if (! fscanf (fp, "%[^=] = ", palName)) {
        return NULL;
    }
    numColors = 0;
    while (fscanf (fp, "0x%x , ", &(colorList[numColors])) == 1) {
        numColors++;
    }

    printf ("Palette '%s' with %d colors:\n", palName, numColors);
    for (i=0; i<numColors; i++) {
        printf ("  %d: 0x%06x\n", i, colorList[i]);
    }

    pal = malloc (sizeof (*pal));
    pal->name = strdup (palName);
    pal->numColors = numColors;
    pal->colorList = malloc (numColors * sizeof (unsigned int));
    pal->shadeList = NULL;
    for (i=0; i<numColors; i++) {
        pal->colorList[i] = colorList[i];
    }

    return pal;
}

int main (int argc, char *argv[]) {

    const int gridSize     = 10;
    const double fieldSize = 0.25;

    const int fColNum      = 5;
    const int fShadeNum    = 256;

    const int animDelay    = 20;

    LedGrid lg;

    FILE *palFile;
    char palName[256];
    int numPals;
    int currPal;
    unsigned int colorList[128];
    int numColors;
    int res;
    double dt = 0.05;
    double gamma = 1.0;

    PaletteType palList[256];
    PaletteType pal;

    //unsigned char *shadeList = NULL;
    int i, j, k;

    int col, row;
    unsigned int color;
    int running = 1;

    pthread_t animationThread;

    void SignalHandler (int arg) {
        running = 0;
        pthread_join (animationThread, NULL);
        endwin ();
        LedGrid_AllOff (lg);
        LedGrid_Free (lg);
        exit (0);
    }

    // getColor --
    //
    //     Retourniert einen 32Bit Integerwert mit dem RGB-Wert der Farbe,
    //     die in der Palette 'pal' den double Wert 't' hat.
    //
    unsigned int getColor (PaletteType pal, double t) {
        assert (pal != NULL);
        assert (t >=0 && t <= 1.0);

        unsigned char red, green, blue;
        int pInd;

        pInd  = (int) (t * fShadeNum * (pal->numColors-1));
        red   = pal->shadeList[pInd*3+0];
        green = pal->shadeList[pInd*3+1];
        blue  = pal->shadeList[pInd*3+2];

        return (red << 16) | (green << 8) | blue;
    }

    // calculateShades --
    //
    //     Berechnet die Zwischenfarben in der Palette 'pal'.
    //
    void calculateShades (PaletteType pal) {
        assert (pal != NULL);

        unsigned char red1, green1, blue1;
        unsigned char red2, green2, blue2;
        unsigned char red,  green,  blue;
        int i, j, k;

        if (pal->shadeList != NULL) {
            free (pal->shadeList);
        }

        pal->shadeList = calloc (fShadeNum*(pal->numColors-1),
                3*sizeof (unsigned char));
        for (i=0; i<pal->numColors-1; i++) {
            for (j=0; j<fShadeNum; j++) {
               red1   = (pal->colorList[i]   >> 16) & 0xff;
               green1 = (pal->colorList[i]   >>  8) & 0xff;
               blue1  = (pal->colorList[i]        ) & 0xff;
               red2   = (pal->colorList[i+1] >> 16) & 0xff;
               green2 = (pal->colorList[i+1] >>  8) & 0xff;
               blue2  = (pal->colorList[i+1]      ) & 0xff;
               red    = red1   + ((double)j/fShadeNum)*(red2-red1);
               green  = green1 + ((double)j/fShadeNum)*(green2-green1);
               blue   = blue1  + ((double)j/fShadeNum)*(blue2-blue1);
               k = 3*(i*fShadeNum + j);
               pal->shadeList[k+0] = red;
               pal->shadeList[k+1] = green;
               pal->shadeList[k+2] = blue;
           }
        }
    }

    // colorFunc0[123] --
    //
    //     Funktionen zur Erzeugung des Plasma-Effektes.
    //
    double colorFunc01 (double x, double y, double t,
            double p1) {
        return sin (x * p1 + t);
    }

    double colorFunc02 (double x, double y, double t,
            double p1, double p2, double p3) {
        return sin (p1*(x*sin (t/p2)+y*cos (t/p3))+t);
    }

    double colorFunc03 (double x, double y, double t,
            double p1, double p2) {
        double cx, cy;

        cx = x + 0.5 * sin (t/p1);
        cy = y + 0.5 * cos (t/p2);
        return sin (sqrt (100.0*(cx*cx + cy*cy)+1.0)+t);
    }

    //-------------------------------------------------------------------------
    //
    // AnimationThreadFunc --
    //
    //     Funktion, die in einem Thread zur Animation der Farben ausgefuehrt
    //     wird.
    //
    // Argument:
    //     arg  Zeiger auf ein LedGrid-Objekt.
    //
    void *AnimationThreadFunc (void *arg) {
        LedGrid lg;
	double x, y, dx, dy;
	double v1, v2, v3, v, cx, cy;
	double time;
        int col, row;

        lg = (LedGrid) arg;
	dx = dy = fieldSize / (gridSize - 1);
	time = 0.0;
	while (running) {
	    y = fieldSize / 2.0;
	    for (row=0; row<gridSize; row++) {
		x = -fieldSize / 2.0;
		for (col=0; col<gridSize; col++) {

                    // Color function 1.
                    //
		    v1 = colorFunc01 (x, y, time, 10.0);

                    // Color function 2.
                    //
                    v2 = colorFunc02 (x, y, time, 10.0, 2.0, 3.0);

                    // Color function 3.
                    //
                    v3 = colorFunc03 (x, y, time, 5.0, 3.0);

		    v = v1 + v2 + v3;
		    v = (v + 3.0) / 6.0;

		    color = getColor (palList[currPal], v);
		    LedGrid_SetColorInt (lg, col, row, color);

		    x += dx;
		}
		y -= dy;
	    }
	    LedGrid_Show (lg);
	    delay (animDelay);
	    time += dt;
	}
    }

    //-------------------------------------------------------------------------
    //
    // (main)
    //
    if (argc != 2) {
        fprintf (stderr, "usage: %s <paletteFile>\n", basename (argv[0]));
        exit (1);
    }
    palFile = fopen (argv[1], "r");
    if (palFile == NULL) {
        fprintf (stderr, "ERROR: couldn't open file '%s'!\n", argv[1]);
        exit (1);
    }
    numPals = 0;
    while (! feof (palFile)) {
/*
        if (! fscanf (palFile, "%[^=] = ", palName)) {
            break;
        }
        numColors = 0;
        while ((res = fscanf (palFile, "0x%x , ",
                    &(colorList[numColors]))) == 1) {
            numColors++;
        }
        printf ("Palette '%s' with %d colors:\n", palName, numColors);
        for (i=0; i<numColors; i++) {
            printf ("  %d: 0x%06x\n", i, colorList[i]);
        }

        pal = malloc (sizeof (PaletteType));
        pal->name = strdup (palName);
        pal->numColors = numColors;
        pal->colorList = malloc (numColors * sizeof (unsigned int));
        pal->shadeList = NULL;
        for (i=0; i<numColors; i++) {
            pal->colorList[i] = colorList[i];
        }
*/
        if ((pal = Pal_Init (palFile)) == NULL) {
            break;
        }

        palList[numPals] = pal;
        printf ("Calculate shades...");
        fflush (stdout);
        calculateShades (pal);
        printf ("done!\n");
        numPals++;
    }
    printf ("Reading completed!\n");
    fclose (palFile);

    currPal = 0;

    signal (SIGINT, SignalHandler);

    initscr ();
    raw ();
    keypad (stdscr, TRUE);
    noecho ();

    lg = LedGrid_Init (gridSize, gridSize, gamma);nstall

    pthread_create (&animationThread, NULL, &AnimationThreadFunc, lg);

    while (running) {
        int ch;

        clear ();
        printw ("------------------------\n");
        printw ("Gamma value    : %1.3f\n", gamma);
        printw ("e: decr; r: incr\n");
        printw ("------------------------\n");
        printw ("Animation dT   : %0.3f\n", dt);
        printw ("a: decr; s: incr\n");
        printw ("------------------------\n");
        printw ("Current palette: %s\n", palList[currPal]->name);
        printw ("q: prev; w: next\n");
        for (i=0; i<palList[currPal]->numColors; i++) {
            printw ("  > Color %d: 0x%06x\n", i,
                    palList[currPal]->colorList[i]);
        }
        printw ("------------------------\n");
        printw ("x: quit\n");
        printw ("------------------------\n");
        refresh ();

        ch = wgetch (stdscr);
        switch (ch) {
            case 'x':
                running = 0;
                break;
            case 'q':
                if (currPal > 0) {
                    currPal--;
                }
                break;
            case 'w':
                if (currPal < numPals-1) {
                    currPal++;
                }
                break;
            case 'a':
                if (dt > 0.01) {
                    dt-=0.01;;
                }
                break;
            case 's':
                dt+=0.01;
                break;
            case 'e':
                if (gamma > 1.0) {
                    gamma -= 0.1;
                }
                LedGrid_SetGamma (lg, gamma);
                break;
            case 'r':
                gamma += 0.1;
                LedGrid_SetGamma (lg, gamma);
                break;
        }
    }

    SignalHandler (0);

    return 0;
}

