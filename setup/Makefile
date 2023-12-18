## Start-up and tear-down targets

.PHONY: dcup
dcup:
	docker-compose -f docker-compose.yaml up -d \
		mex-metadata-svc  mex-query-svc  mex-auth-svc  mex-webapp  mex-config-svc \
		mex-proxy  pairgen \
		mex-db  mex-redis  mex-solr-1  mex-solr-2  mex-solr-3  mex-zoo-1  mex-zoo-2  mex-zoo-3  mex-redis-commander

.PHONY: dcdownv
dcdownv:
	docker-compose down --volumes --remove-orphans

## Demo scenario targets

.PHONY: demo-load-mesh
demo-load-mesh:
	./tools/mexctl/mexctl.sh --config @./config.yaml.example --instance demo \
		blobs create -n mesh-blob -t sqlitez -z -f ./scenarios/demo/codings/minimesh.db

.PHONY: demo-load-config
demo-load-config:
	./tools/mexctl/mexctl.sh --config @./config.yaml.example --instance demo \
		config send --dir ./scenarios/demo/config

.PHONY: demo-load-items
demo-load-items:
	./tools/mexctl/mexctl.sh --config @./config.yaml.example --instance demo \
		data load --bulk \
		-i ./scenarios/demo/items/random.items.yaml \
		-i ./scenarios/demo/items/additional.items.yaml \
		-i ./scenarios/demo/items/orgunits.items.yaml

## Croissant scenario targets

.PHONY: croissant-generate-config
croissant-generate-config:
	./tools/mexctl/mexctl.sh --config @./config.yaml.example --instance demo \
		croissant generate \
			--in ./scenarios/croissant/datasets \
			--config-out ./scenarios/croissant/config/test \
			--items-out ./scenarios/croissant/items

.PHONY: croissant-load-config
croissant-load-config:
	./tools/mexctl/mexctl.sh --config @./config.yaml.example --instance demo \
		config send --dir ./scenarios/croissant/config

.PHONY: croissant-load-items
croissant-load-items:
	./tools/mexctl/mexctl.sh --config @./config.yaml.example --instance demo \
		data load --bulk \
		-i ./scenarios/croissant/items/items.dataset.yaml \
		-i ./scenarios/croissant/items/items.recordset.yaml

## Generic targets

.PHONY: show-tree
show-tree:
	./tools/mexctl/mexctl.sh --config @./config.yaml.example --instance demo \
		tree show -n OrganizationalUnit -l parentDepartment -d label

.PHONY: delete-items
delete-items:
	./tools/mexctl/mexctl.sh --config @./config.yaml.example --instance demo \
		items delete
