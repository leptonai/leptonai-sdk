import { FC } from "react";
import { Card } from "@lepton-dashboard/components/card";
import { NewBreadcrumbProps } from "antd/es/breadcrumb/Breadcrumb";
import { Breadcrumb } from "antd";
import { css } from "@emotion/react";

export const BreadcrumbHeader: FC<NewBreadcrumbProps> = ({ items }) => {
  return (
    <Card
      css={css`
        padding: 8px 16px;
      `}
    >
      <Breadcrumb items={items} />
    </Card>
  );
};
