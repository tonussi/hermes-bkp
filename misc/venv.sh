#!/usr/bin/env sh
export VENV=venv 
python3 -m venv $VENV
echo "source $VENV/bin/activate"
echo "pip install -r plotters/requirements.txt"
