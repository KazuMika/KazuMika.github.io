#!/bin/bash

DATE=$(date +%Y.%m.%d_%H.%M.%S)
HERE=$(pwd)
LOG_DIR="${HERE}"/.log

if [ ! -d "${LOG_DIR}" ]; then
    mkdir "${LOG_DIR}";
fi






export PYTHONPATH=${PYTHONPATH}:${HERE}/src
export PYTHONPATH=${PYTHONPATH}:${HERE}/src/utils
export PYTHONPATH=${PYTHONPATH}:${HERE}/src/pycalib
export PYENV_ROOT=$HOME/.pyenv
export PATH=$PYENV_ROOT/bin:$PATH

VENV_NAME="sde"

if [ -d "${PYENV_ROOT}" ]; then
    pyenv local "${VENV_NAME}"
fi

if command -v pyenv 1>/dev/null 2>&1; then
    eval "$(pyenv init -)"
fi
eval "$(pyenv virtualenv-init -)"


data_dir="${HERE}/data"
results_dir="${HERE}/results"

which=${1:-run}


THIS_START=$(date +%Y.%m.%d_%H.%M.%S)
echo "[START: $which] [$THIS_START] ##########################################"


if [ "${which}" == "edit" ]; then
    _data_dir="${data_dir}/"
    pl="00"
    pid="00"
    start=32
    end=32
    frames=("32" "162" "442")

    ext='npz'
    L=3
    numTrial=40
    thErrReProj=15.0
    rateForThOKNum=0.55

    # frame.隠れてない個数.id
    # 32.4.10
    # 162.3.5
    # 32.2.0
    # 442.2.6


    for frm in $(seq $start $end); do
    #for frm in "${frames[@]}" ; do
        python ./src/main.py \
            --which "$which" \
            --data_dir "${_data_dir}" \
            --L "${L}" \
            --numTrial "${numTrial}" \
            --rateForThOKNum "${rateForThOKNum}" \
            --thErrReProj "${thErrReProj}" \
            --pl "${pl}" \
            --pid "${pid}" \
            --frm_num ${frm} \
            --ext "${ext}" 2>&1 | tee "$LOG_DIR/$which.log.$DATE"
    done

fi




if [ "${which}" == "make_grid" ]; then

    ext='mp4'
    _data_dir="${data_dir}/pl00/pl00_id00/"
    _results_dir="${results_dir}/grid"

    python ./src/main.py \
        --which "$which" \
        --data_dir "${_data_dir}" \
        --results_dir "${_results_dir}" \
        --ext "${ext}" 2>&1 | tee "$LOG_DIR/$which.log.$DATE"
fi


if [ "${which}" == "setup_venv" ]; then
    git clone https://github.com/pyenv/pyenv.git "${HOME}"/.pyenv
    git clone https://github.com/yyuu/pyenv-virtualenv.git "${HOME}"/.pyenv/plugins/pyenv-virtualenv
    pyenv install -f 3.9.14
    pyenv virtualenv 3.9.14 "${VENV_NAME}"
    pyenv local "${VENV_NAME}"
    pip install --upgrade pip
    pip install -r requirements.txt
fi


if [ "${which}" == imgrid ]; then
    dir_names=("00" "01" "02" "03")
    mkdir -p "${HERE}/temp2/"

    for dir_name in "${dir_names[@]}"; do
        im_lst="${HERE}/results/hand_edit/cam${dir_name}/"
        grid="4x4"
        im_save_fn="${HERE}/temp2/cam${dir_name}"
        im_save_fn="${HERE}/temp2/cam${dir_name}"
        mkdir -p  "${im_save_fn}"
        im_save_fn=$im_save_fn/results.png
        imgrid --grid "${grid}" --im_lst "${im_lst}" --im_save_fn "${im_save_fn}"
    done
fi



if [ "${which}" == git ]; then
    comment=$2
    git add .
    git commit -m "${DATE} $comment"
    # git push
fi


if [ "${which}" == license ]; then
    pip-licenses --with-system --order=license > docs/LICENSE.txt
fi


if [ "${which}" == rsync ]; then
    # rsync -auvz "${HERE}/" "jp:${HOME}/"
    rsync -avz --delete "${HERE}" "jp:project/"
fi


THIS_END=$(date +%Y.%m.%d_%H.%M.%S)
echo "[END  : $which] [$THIS_END] ############################################"
