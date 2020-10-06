#!/bin/bash


if [[ -z "$KEPTNPROJECT" ]]; then
  KEPTNPROJECT=pacproject
fi
if [[ -z "$KEPTNSERVICE" ]]; then
  KEPTNSERVICE=pacservice
fi
if [[ -z "$KEPTNSTAGE" ]]; then
  KEPTNSTAGE=qualitygate
fi
if [[ -z "$K3SKUBECTL" ]]; then
  K3SKUBECTL=k3s kubectl
fi

echo "Assumes Keptn CLI is configured and points to a Keptn Installation"
read -rsp $'Press ctrl-c to abort. Press any key to continue...\n' -n1 key

echo "-----------------------------------------------"
echo "Step 2 - Install PAC SLI Provider"
"${K3SKUBECTL[@]}" -n keptn apply -f https://raw.githubusercontent.com/grabnerandi/pac-sliprovider/master/deploy/service.yaml

echo "-----------------------------------------------"
echo "Step 3 - Create a Keptn Project for PAC"
wget https://raw.githubusercontent.com/grabnerandi/pac-sliprovider/master/keptnproject/shipyard.yaml
keptn create project "${KEPTNPROJECT}" -s=shipyard.yaml

echo "-----------------------------------------------"
echo "Step 4 - Configure PAC Provider for our Project"
"${K3SKUBECTL[@]}" -n keptn apply -f https://raw.githubusercontent.com/grabnerandi/pac-sliprovider/master/keptnproject/lighthouse-configmap.yaml

echo "-----------------------------------------------"
echo "Step 5 - Create a service"
keptn create service "${KEPTNSERVICE}" -p="${KEPTNPROJECT}"

echo "-----------------------------------------------"
echo "Step 6 - Uploading SLO.yaml"
wget https://raw.githubusercontent.com/grabnerandi/pac-sliprovider/master/keptnproject/slo.yaml
wget https://raw.githubusercontent.com/grabnerandi/pac-sliprovider/master/keptnproject/pac-sliprovider/sli.yaml
keptn add-resource --project="${KEPTNPROJECT}" --service="${KEPTNSERVICE}" --stage="${KEPTNSTAGE}" --resource=slo.yaml
keptn add-resource --project="${KEPTNPROJECT}" --service="${KEPTNSERVICE}" --stage="${KEPTNSTAGE}" --resource=sli.yaml --resourceUri=pac-sliprovider/sli.yaml

echo "-----------------------------------------------"
echo "Step 7 - Executing Quality Gates"
keptn send event start-evaluation --project="${KEPTNPROJECT}" --service="${KEPTNSERVICE}" --stage="${KEPTNSTAGE}" --labels=pacId=FirstPAC,buildId=1
keptn send event start-evaluation --project="${KEPTNPROJECT}" --service="${KEPTNSERVICE}" --stage="${KEPTNSTAGE}" --labels=pacId=MountainPAC,buildId=2
keptn send event start-evaluation --project="${KEPTNPROJECT}" --service="${KEPTNSERVICE}" --stage="${KEPTNSTAGE}" --labels=pacId=FuturePAC,buildId=3
keptn send event start-evaluation --project="${KEPTNPROJECT}" --service="${KEPTNSERVICE}" --stage="${KEPTNSTAGE}" --labels=pacId=JurassicPAC,buildId=4


echo "-----------------------------------------------"
echo "DONE!!"