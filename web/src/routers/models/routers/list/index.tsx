import { FC, useMemo, useState } from "react";
import { Col, Input, Row } from "antd";
import { SearchOutlined } from "@ant-design/icons";
import { useInject } from "@lepton-libs/di";
import { ModelService } from "@lepton-dashboard/services/model.service.ts";
import { useStateFromObservable } from "@lepton-libs/hooks/use-state-from-observable.ts";
import { ModelGroupCard } from "../../../../components/model-group-card";
import { DeploymentService } from "@lepton-dashboard/services/deployment.service.ts";

export const List: FC = () => {
  const modelService = useInject(ModelService);
  const groupedModels = useStateFromObservable(() => modelService.groups(), []);
  const [search, setSearch] = useState<string>("");
  const filteredModels = useMemo(() => {
    return groupedModels.filter(
      (e) => JSON.stringify(e).indexOf(search) !== -1
    );
  }, [groupedModels, search]);
  const deploymentService = useInject(DeploymentService);
  const deployments = useStateFromObservable(
    () => deploymentService.list(),
    []
  );
  return (
    <Row gutter={[8, 24]}>
      <Col flex={1}>
        <Input
          autoFocus
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          prefix={<SearchOutlined />}
          placeholder="Search"
        />
      </Col>
      <Col span={24}>
        <Row gutter={[16, 16]} wrap>
          {filteredModels.map((group) => (
            <Col flex="1" key={group.name}>
              <ModelGroupCard
                deploymentCount={
                  deployments.filter((i) =>
                    group.data.some((m) => m.id === i.photon_id)
                  ).length
                }
                group={group}
              />
            </Col>
          ))}
        </Row>
      </Col>
    </Row>
  );
};
