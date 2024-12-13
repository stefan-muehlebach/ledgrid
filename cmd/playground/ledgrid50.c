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
#include "PiPack.h"

/*-----------------------------------------------------------------------------
 *
 * (main)
 *
 *     Versuch, die Plasma-Funktion von PixelController zu imitieren.
 *
 */

int main (int argc, char *argv[]) {
    LedGrid lg;
    const int gridSize = 10;
    int x = 0, y = 0;
    int running = 1;

    int frameCount = 1;
    int timeDisplacement;
    float xc;
    float calculation1, calculation2;
    int aaa;
    float yc, s1, s2, s3, s;
    int aa;
    int i;

    int colorSetSize = 5;
    int colorSet[5] = {0x492d61, 0x048091, 0x61c155, 0xf2d43f, 0xd1026c};

/*
    int colorSetSize = 3;
    int colorSet[3] = {0x000000, 0x00ffff, 0x00ff00};
*/

    int colorArray[256];

/*
    unsigned char redVals[6]   = {0x00, 0x00, 0x00, 0x00, 0x80, 0x80};
    unsigned char greenVals[6] = {0x00, 0xff, 0x00, 0x80, 0xff, 0x80};
    unsigned char blueVals[6]  = {0x00, 0xff, 0x80, 0x80, 0xff, 0xff};
    // unsigned char redArray[255], greenArray[255], blueArray[255];
    unsigned char colorArray[256];
*/

/*
    unsigned char colorArray[255] = {
        0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1,
        0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1,
        0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1,
        0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1,
        0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1,
        0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1,
        0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1,
        0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1,
        0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1,
        0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1,
        0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1,
        0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1,
        0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6
    };
*/

    int calcSmoothColor(int col1, int col2, int pos) {
        int b = col1 & 0xff;
        int g = col1 >> 8 & 0xff;
        int r = col1 >> 16 & 0xff;
        int b2 = col2 & 0xff;
        int g2 = col2 >> 8 & 0xff;
        int r2 = col2 >> 16 & 0xff;

        int mul = pos * colorSetSize;
        int oppositeColor = 255 - mul;

        r = r * mul + r2 * oppositeColor >> 8;
        g = g * mul + g2 * oppositeColor >> 8;
        b = b * mul + b2 * oppositeColor >> 8;

        return r << 16 | g << 8 | b;
    }

    float boarderCount = 255.0 / (float) colorSetSize;
    //printf ("boarderCount: %f\n", boarderCount);
    for (i=0; i<256; i++) {
        int ofs = 0;
        int pos = i;
        while (pos > boarderCount) {
            pos -= boarderCount;
            //pos = (int)(pos - boarderCount);
            ofs++;
        }
        int targetOfs = ofs + 1;
        //printf ("[%03d] ofs: %d, targetOfs: %d\n", i, ofs, targetOfs);
        colorArray[i] = calcSmoothColor(colorSet[(targetOfs % colorSetSize)], \
                colorSet[(ofs % colorSetSize)], pos);
    }

    //for (i=0; i<256; i++) {
    //    printf ("%3d: 0x%06x\n", i, colorArray[i]);
    //}

    void SignalHandler (int arg) {
        endwin ();
        LedGrid_AllOff (lg);
        LedGrid_Free (lg);
        exit (0);
    }

    float radians (float degrees) {
        return degrees * 0.01745329;
    }

    signal (SIGINT, SignalHandler);
    lg = LedGrid_Init (gridSize, gridSize, 1.5);

    while (running) {
        xc = 20.0;
        timeDisplacement = frameCount++;

        calculation1 = sin (radians (timeDisplacement * 0.6165562));
        calculation2 = sin (radians (timeDisplacement * -3.635226));

        aaa = 128;
        for (x=0; x<gridSize; x++, xc++) {
            yc = 20.0;
            s1 = aaa + aaa * sin (radians (xc) * calculation1);
            for (y=0; y<gridSize; y++, yc++) {
                s2 = aaa + aaa * sin (radians (yc) * calculation2);
                s3 = aaa + aaa * sin (radians ((xc + yc + timeDisplacement * 3.0) / 2.0));
                s  = (s1 + s2 + s3) / 255.0;
                aa = (int)(s * 255.0 + 0.5);
                LedGrid_SetColorInt (lg, x, y, colorArray[aa % 255]);
            }
        }
        LedGrid_Show (lg);
        delay (50);
    }

    SignalHandler (0);

    return 0;
}

