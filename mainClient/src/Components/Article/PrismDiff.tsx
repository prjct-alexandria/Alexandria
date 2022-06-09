import React from "react";
import PropTypes from "prop-types";
const Prism = require('prismjs')
const diff = require('diff');

const PrismDiff = ({ sourceContent = "", targetContent = ""}) => {
    let groups = diff.diffWords(sourceContent, targetContent);

    const mappedNodes = groups.map((group: { value: any; added: any; removed: any; }, i: number) => {
        const { value, added, removed } = group;
        let nodeStyles;
        if (added) nodeStyles = "added";
        if (removed) nodeStyles = "removed";

        // Using dangerouslySetInnerHTML with the Node rendering API
        // Note: is dangerous
        return (
            <span key={i}
                className={nodeStyles}
                dangerouslySetInnerHTML={{
                    __html: Prism.highlight(
                        value,
                        Prism.languages.javascript,
                        "javascript"
                    ),
                }}
            />
        );
    });

    return <span>{mappedNodes}</span>;
};

PrismDiff.propTypes = {
    sourceContent: PropTypes.string,
    targetContent: PropTypes.string,
};

export default PrismDiff;