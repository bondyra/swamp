

import React, { memo } from 'react';
import SingleFieldPicker from './pickers/SingleFieldPicker';


export default memo(({ op, change }) => {
    const options = [
        {value: "==", description: "Left equals to right"},
        {value: "!=", description: "Left not equals to right"},
        {value: "like", description: "Left like right"},
        {value: "not like", description: "Left not like right"},
        {value: "contains", description: "Left contains right"},
        {value: "not contains", description: "Left not contains right"},
    ]
  return (
    <SingleFieldPicker value={op} updateData={(v) => {change(v)}} options={options}
    valuePlaceholder="op" popperPrompt="Select operation"/>
  );
});
