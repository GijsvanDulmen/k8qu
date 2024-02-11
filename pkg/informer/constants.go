package informer

import logger "k8qu/pkg/log"

const finalizerName = "k8qu.io"

var log = logger.Logger()

type JobReconcileRequest func(forQueue string)
