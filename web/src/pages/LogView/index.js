
import React from 'react';
import { Segment, Tab } from 'semantic-ui-react';
import LogViewTable from '../../components/LogViewTable';


const LogView = () => {
 
    let panes = [
      {
        menuItem: '错误日志',
        render: () => (
          <Tab.Pane attached={true}>
                <LogViewTable type="error" />
          </Tab.Pane>
        ),
      },
       {
        menuItem: '常规日志',
        render: () => (
          <Tab.Pane attached={false}>
            <LogViewTable type="common" />
          </Tab.Pane>
        ),
      },
    ];  
  return (
    <Segment>
      <Tab menu={{ secondary: true, pointing: true }} panes={panes} />
    </Segment>
  );
};

export default LogView;
