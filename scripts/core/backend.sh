#!/usr/bin/env bash

source ./common.sh

#setup fh-aaa, fh-scm, fh-messaging, fh-metrics, gitlab-shell, ups, supercore (supercore needs to be running before millicore)

##TODO ups, fh-messaging, fh-metrics
## add these to params once all ready UPS_IMAGE_VERSION="$VERSION_UPS"
## add these once messaging metrics mongo changes in FH_MESSAGING_IMAGE_VERSION="$VERSION_MESSAGING",FH_METRICS_IMAGE_VERSION="$VERSION_METRICS"
oc new-app --param=FH_AAA_IMAGE_VERSION="$VERSION_AAA",FH_SUPERCORE_IMAGE_VERSION="$VERSION_SUPERCORE",FH_SCM_IMAGE_VERSION="$VERSION_SCM",GITLAB_SHELL_IMAGE_VERSION="$VERSION_GITLAB",MONGODB_FHSUPERCORE_PASSWORD="$MONGO_PASSWORD",REDIS_PASSWORD="$REDIS_PASSWORD",FH_MESSAGING_API_KEY="$MESSAGING_API_KEY",FH_METRICS_API_KEY="$METRICS_API_KEY",FH_MILLICORE_API_KEY="$MILLICORE_API_KEY",FH_MILLICORE_SERVICE_KEY="$MILLICORE_SERVICE_KEY",MONGODB_FHAAA_PASSWORD="$MONGO_PASSWORD" -f "${TEMPLATES_DIR}/generated/fh-core-backend.json"

#gitlab-shell
#these need to be done as the admin user (root)
#oc create -f "${TEMPLATES_DIR}/gitlab-shell/scc-anyuid-with-chroot.json"

#oc adm policy add-scc-to-user anyuid-with-chroot system:serviceaccount:core:default