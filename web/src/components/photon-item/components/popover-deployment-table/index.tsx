import { FC, useMemo } from "react";
import { Popover, Table, Tag } from "antd";
import { Deployment } from "@lepton-dashboard/interfaces/deployment.ts";
import { Photon } from "@lepton-dashboard/interfaces/photon.ts";
import { DeploymentStatus } from "@lepton-dashboard/components/deployment-status";
import { useNavigate } from "react-router-dom";
import { css as classNameCss } from "@emotion/css";
import { Link } from "@lepton-dashboard/components/link";
import { DeploymentIcon } from "@lepton-dashboard/components/icons";
import { Description } from "@lepton-dashboard/components/description";
import { DateParser } from "@lepton-dashboard/components/date-parser";
import { css } from "@emotion/react";

export const PopoverDeploymentTable: FC<{
  photon: Photon;
  deployments: Deployment[];
}> = ({ photon, deployments }) => {
  const navigate = useNavigate();
  const color = useMemo(() => {
    const running = deployments.some((d) => d.status.state === "Running");
    const hasDeployments = deployments.length > 0;
    if (running) {
      return "success";
    } else if (hasDeployments) {
      return "processing";
    } else {
      return "default";
    }
  }, [deployments]);
  return (
    <Tag
      color={color}
      css={css`
        margin-right: 0;
      `}
    >
      <Description.Item
        icon={<DeploymentIcon />}
        description={
          <Popover
            open={deployments.length > 0 ? undefined : false}
            placement="bottomLeft"
            content={
              <Table
                rowClassName={classNameCss`cursor: pointer;`}
                size="small"
                showHeader={false}
                pagination={false}
                bordered
                rowKey="id"
                onRow={(record) => {
                  return {
                    onClick: () => navigate(`/deployments/detail/${record.id}`),
                  };
                }}
                columns={[
                  {
                    title: "Status",
                    dataIndex: ["status", "state"],
                    render: (state) => <DeploymentStatus status={state} />,
                  },
                  {
                    title: "Name",
                    dataIndex: "name",
                    render: (v) => v,
                  },
                  {
                    title: "Created",
                    dataIndex: "created_at",
                    render: (data) => <DateParser detail date={data} />,
                  },
                ]}
                dataSource={deployments}
              />
            }
          >
            <span>
              <Link to={`/deployments/list/${photon.name}`} relative="route">
                {deployments.length > 0 ? deployments.length : "No"}{" "}
                {deployments.length > 1 ? "deployments" : "deployment"}
              </Link>
            </span>
          </Popover>
        }
      />
    </Tag>
  );
};
