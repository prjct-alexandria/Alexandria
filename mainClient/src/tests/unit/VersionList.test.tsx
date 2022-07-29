import * as React from "react";
import VersionList from "../../Components/Article/VersionList";
import {render, screen, waitFor} from "@testing-library/react";
import {MemoryRouter} from "react-router-dom";
import MRList from "../../Components/Article/MRList";

// Arrange: Making a fake response
const versionList = [
    {
        "articleID": 1,
        "versionID": 1,
        "title":"Web-Development research",
        "owners":["Jane Doe", "John Doe"],
        "content":"\nLorem ipsum dolor sit amet, consectetur adipiscing elit. Duis commodo felis sed dui suscipit faucibus. Duis eget imperdiet augue, mollis imperdiet urna. Donec ac elit sodales, interdum enim eget, tristique metus. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. In vehicula eu ligula posuere dapibus. Curabitur lorem ligula, ornare a purus vitae, blandit ullamcorper sem. Nulla dignissim egestas luctus. Sed volutpat consequat velit, non bibendum neque aliquet non. Mauris malesuada hendrerit eros ut interdum. Aliquam erat volutpat. Fusce nec rhoncus erat. Cras auctor elit in aliquet suscipit. In sed rutrum leo. Mauris tincidunt ante id leo facilisis congue.",
        "status": "pending"
    },
    {
        "articleID": 1,
        "versionID": 2,
        "title":"Web-Development research v2",
        "owners":["Jane Doe", "John Doe"],
        "content":"\nLorem ipsum dolor sit amet, adipiscing elit. Duist commodos felis sed dui suscipit faucibus. Duis eget imperdiet augue, mollis imperdiet urna. Donec ac elit sodales, interdum enim eget, tristique metus. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. In vehicula eu ligula posuere dapibus. Curabitur lorem ligula, ornare a purus vitae, blandit ullamcorper sem. Nulla dignissim egestas luctus. Sed volutpat consequat velit, non bibendum neque aliquet non. Mauris malesuada hendrerit eros ut interdum. Aliquam erat volutpat. Fusce nec rhoncus erat. Cras auctor elit in aliquet suscipit. In sed rutrum leo. Mauris tincidunt leo facilisis congue.",
        "status": "accepted"
    },
    {
        "articleID": 1,
        "versionID": 3,
        "title":"Third version",
        "owners":["Jane Doe", "John Doe"],
        "content":"some content",
        "status": "rejected"
    }
]

let mockedFetch = null;

describe('VersionList', () => {
    // Arrange: Mock the fetch
    beforeEach(() => {
        const response = {
            status: 200,
            ok: true,
            json: jest.fn().mockResolvedValue(versionList) };
        mockedFetch = jest.fn().mockResolvedValue(response);
        global.fetch = mockedFetch;
    })

    // Test whether the first article is a part of the articleList
    test("test versions", async () => {

        // Act: Render the article list
        render(
            <MemoryRouter>
                <VersionList />
            </MemoryRouter>
        );

        // Wait until it is finished loading
        const loadingSpinner = await screen.findByTestId("loadingSpinner")
        await waitFor(() => {
            expect(loadingSpinner).not.toBeInTheDocument()
        })

        // Assert that the fetch was used
        expect(mockedFetch).toHaveBeenCalledTimes(2)

        // Perform assertions for each article
        versionList.forEach((versionListElement) => {
            const i = versionListElement.versionID
            expect(screen.getByTestId("title" + i)).toHaveTextContent(versionListElement.title)
            versionListElement.owners.forEach((owner)=> {
                expect(screen.getByTestId("owners" + i)).toHaveTextContent(owner)
            })
            expect(screen.getByTestId("status" + i)).toHaveTextContent(versionListElement.status)
        })
    })
})

let alert = {
    message: "message"
}
describe('VersionList error', () => {
    // Arrange: Mock the fetch
    beforeEach(() => {
        const response = {
            status: 404,
            ok: false,
            json: jest.fn().mockResolvedValue(alert) };
        mockedFetch = jest.fn().mockResolvedValue(response);
        global.fetch = mockedFetch;
    })

    // Test whether the first article is a part of the articleList
    test("test versions", async () => {
        // Act: Render the article list
        render(
            <MemoryRouter>
                <VersionList />
            </MemoryRouter>
        );

        // Wait until it is finished loading
        const loadingSpinner = await screen.findByTestId("loadingSpinner")
        await waitFor(() => {
            expect(loadingSpinner).not.toBeInTheDocument()
        })

        // Assert that the fetch was used
        expect(mockedFetch).toHaveBeenCalledTimes(2)

        // Perform assertions for error message
        expect(screen.getByRole("alert")).toHaveTextContent(
            "Error: Something went wrong. Error: " + alert.message
        )
    })
})