BIN_DIR:=bin
WWW_DIR:=www
SSL_DIR:=${BIN_DIR}/ssl
BUILD_DIR:=build

TARGET:=${BIN_DIR}/server
SERVER_CONF:=${BIN_DIR}/survana.json

MAKEFILE_PATH:= $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR:= $(dir $(MAKEFILE_PATH))

COVER:=$(strip $(shell go tool | grep cover))
ifdef COVER
    COVER:=-cover
endif

COVERFILE := survana.coverage

#OSX
OSX_PROJECT_NAME:=Survana
OSX_TARGET:=bin/${OSX_PROJECT_NAME}.app
OSX_BUILD_ARCHIVE:=build/${OSX_PROJECT_NAME}.xcarchive
OSX_BUILD_RESOURCES:=${OSX_BUILD_ARCHIVE}/Products/Applications/Survana.app/Contents/Resources/

OSX_ROOT:=src/platform/osx
OSX_PROJECT_ROOT:=${OSX_ROOT}/${OSX_PROJECT_NAME}
OSX_PROJECT:=${OSX_PROJECT_ROOT}/${OSX_PROJECT_NAME}.xcodeproj

XCODE:=xcodebuild

all: ${TARGET}

${TARGET}:
	GOPATH=${CURRENT_DIR} go get code.google.com/p/goauth2/oauth
	GOPATH=${CURRENT_DIR} go get labix.org/v2/mgo
	GOPATH=${CURRENT_DIR} go install server

test:
	go test ${COVER} neuroinformatics.harvard.edu/survana

clean:
	@rm -f ${TARGET}
	
osx: ${TARGET} ${OSX_TARGET}

${OSX_TARGET}: ${OSX_BUILD_ARCHIVE} ${OSX_PROJECT}
	${XCODE} -project ${OSX_PROJECT} -exportArchive -exportFormat APP -archivePath ${OSX_BUILD_ARCHIVE} -exportPath ${OSX_TARGET}
	
${OSX_BUILD_ARCHIVE}: ${OSX_PROJECT}
	${XCODE} -project ${OSX_PROJECT} -scheme ${OSX_PROJECT_NAME} archive -archivePath ${OSX_BUILD_ARCHIVE}
	@mkdir -p ${OSX_BUILD_RESOURCES}/services/server
	@cp -r ${BIN_DIR}/server ${BIN_DIR}/ssl ${OSX_BUILD_RESOURCES}/services/server/

osx-clean:
	rm -rf ${OSX_TARGET} ${OSX_BUILD_ARCHIVE}

cover: ${TARGET}
	GOPATH=${CURRENT_DIR} go test -coverprofile=${COVERFILE} -covermode=count neuroinformatics.harvard.edu/survana
	GOPATH=${CURRENT_DIR} go tool cover -html=${COVERFILE}

bench: ${TARGET}
	GOPATH=${CURRENT_DIR} go test -bench . neuroinformatics.harvard.edu/survana
