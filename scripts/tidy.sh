CURRENT_DIR=$(pwd)
SEARCH_DIR=$CURRENT_DIR
SEARCH_TOTAL=0
EXECUTE_CMD="go mod tidy"

function GoTidy() {
	for file in $(ls $1); do
		local target="$1/$file"
		if [ -d $target ]; then
			cd $target
			#case
			if [ -f "go.mod" ];then
				`$EXECUTE_CMD`
				if [ $? -ne 0 ]; then
					break
				fi
				echo "process $target"
				let "SEARCH_TOTAL+=1"
			fi
			GoTidy $target
		fi
	done
}

#test
if [ $# -ne 0 ]; then
	if [ -d $1 ]; then
		cd $1
		SEARCH_DIR=$(pwd)
	else
		echo "$1 is not exist directory."
		exit
	fi
fi
#start
GoTidy $SEARCH_DIR
echo "processed($SEARCH_TOTAL)."
#back directory
cd $CURRENT_DIR
