import { FC, useEffect } from "react";
import { Navigate, Route, Routes } from "react-router-dom";
import { List } from "@lepton-dashboard/routers/photons/routers/list";
import styled from "@emotion/styled";
import { Detail } from "@lepton-dashboard/routers/photons/routers/detail";
import { Versions } from "@lepton-dashboard/routers/photons/routers/versions";
import { useInject } from "@lepton-libs/di";
import { TitleService } from "@lepton-dashboard/services/title.service";
const Container = styled.div`
  flex: 1 1 auto;
`;
export const Photons: FC = () => {
  const titleService = useInject(TitleService);
  useEffect(() => {
    titleService.setTitle("Photons");
  }, [titleService]);
  return (
    <Container>
      <Routes>
        <Route path="list" element={<List />} />
        <Route path="versions/:name" element={<Versions />} />
        <Route path="detail/:id" element={<Detail />} />
        <Route path="*" element={<Navigate to="list" replace />} />
      </Routes>
    </Container>
  );
};
