#!/bin/sh

for i in jack_*.png; do 
    newname=$(echo $i | sed "s/jack_/J_/")
    echo $newname
    mv $i $newname
done

for i in king_*.png; do 
    newname=$(echo $i | sed "s/king_/K_/")
    echo $newname
    mv $i $newname
done

for i in queen_*.png; do 
    newname=$(echo $i | sed "s/queen_/Q_/")
    echo $newname
    mv $i $newname
done

for i in ace_*.png; do 
    newname=$(echo $i | sed "s/ace_/A_/")
    echo $newname
    mv $i $newname
done
