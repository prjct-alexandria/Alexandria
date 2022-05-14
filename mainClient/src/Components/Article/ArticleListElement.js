import React from 'react';
import '../../App.js';
import {Link} from "react-router-dom";

export default class Article extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            open: false,
            buttonText: '▼'
        };
    }

    handleButtonClick = () => {
        this.setState((state) => {
            let txt = '▼';
            if (!state.open) {
                txt = '▲'
            }

            return {
                open: !state.open,
                buttonText: txt
            }
        })
    }

    render() {
        return (

            <div className='row row-no-gutters' style={{border: '3px black solid', margin:'3px', textAlign:'left'}}>
                {/*Title and date*/}
                <div className='row col-md-11'>
                    <div className='col-md-9'>
                        <Link to={"/papers/" + this.props.article.id}>
                            {this.props.article.title}
                        </Link>
                    </div>
                    <div className='col-md-2'>
                        {this.props.article.date_created}
                    </div>
                </div>

                {/*Dropdown button*/}
                <div className='col-md-1'>
                    <button type={"button"} className={"button"} onClick={this.handleButtonClick}>{this.state.buttonText}</button>
                </div>

                {/*Dropdown*/}
                {this.state.open && (
                    <div className={'col-md-12'} style={{border:'black solid 1px'}}>
                        <div className='row' style={{fontWeight:'bold'}}>
                            <div className='col-md-6'>
                                Author: {this.props.article.author}<br/>
                                Rating
                            </div>
                        </div>

                        <br/><br/>
                        {this.props.article.description}
                    </div>
                )}
            </div>

        )
    }
}