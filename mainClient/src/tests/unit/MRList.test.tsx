import * as React from "react";
import MRList from "../../Components/Article/MRList";
import {render, screen, waitFor} from "@testing-library/react";
import {MemoryRouter} from "react-router-dom";
import ArticleList from "../../Components/Article/ArticleList";

// Arrange: Making a fake response

const mrList = [
    {
        "sourceTitle": "Source 1",
        "targetTitle": "Target 1",
        "request": {
            requestID: 1,
            articleID: 2,
            sourceVersionID: 3,
            sourceHistoryID: 4,
            targetVersionID: 5,
            targetHistoryID: 6,
            status: "pending",
            conflicted: false
        }
    },
    {
        "sourceTitle": "Source 2",
        "targetTitle": "Target 2",
        "request": {
            requestID: 7,
            articleID: 8,
            sourceVersionID: 9,
            sourceHistoryID: 10,
            targetVersionID: 11,
            targetHistoryID: 12,
            status: "accepted",
            conflicted: true
        }
    },
    {
        "sourceTitle": "Source 3",
        "targetTitle": "Target 3",
        "request": {
            requestID: 13,
            articleID: 14,
            sourceVersionID: 15,
            sourceHistoryID: 16,
            targetVersionID: 17,
            targetHistoryID: 18,
            status: "rejected",
            conflicted: false
        }
    }
]

let mockedFetch = null;

describe('MRList', () => {
    // Arrange: Mock the fetch
    beforeEach(() => {
        const response = {
            status: 200,
            ok: true,
            json: jest.fn().mockResolvedValue(mrList) };
        mockedFetch = jest.fn().mockResolvedValue(response);
        global.fetch = mockedFetch;
    })

    // Test whether the first article is a part of the articleList
    test("test MRs", async () => {

        // Act: Render the article list
        render(
            <MemoryRouter>
                <MRList />
            </MemoryRouter>
        );

        // Wait until it is finished loading
        const loadingSpinner = await screen.findByTestId("loadingSpinner")
        await waitFor(() => {
            expect(loadingSpinner).not.toBeInTheDocument()
        })

        // Assert that the fetch was used
        expect(mockedFetch).toHaveBeenCalledTimes(1)

        // Perform assertions for each article
        mrList.forEach((mrListElement) => {
            const i = mrListElement.request.requestID
            expect(screen.getByTestId("sourceTitle" + i)).toHaveTextContent(mrListElement.sourceTitle)
            expect(screen.getByTestId("targetTitle" + i)).toHaveTextContent(mrListElement.targetTitle)
            expect(screen.getByTestId("status" + i)).toHaveTextContent(mrListElement.request.status)
        })
    })
})

let alert = {
    message: "message"
}
describe('MRList error', () => {
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
    test("test requests", async () => {
        // Act: Render the article list
        render(
            <MemoryRouter>
                <MRList />
            </MemoryRouter>
        );

        // Wait until it is finished loading
        const loadingSpinner = await screen.findByTestId("loadingSpinner")
        await waitFor(() => {
            expect(loadingSpinner).not.toBeInTheDocument()
        })

        // Assert that the fetch was used
        expect(mockedFetch).toHaveBeenCalledTimes(1)

        // Perform assertions for error message
        expect(screen.getByRole("alert")).toHaveTextContent(
            "Error: Something went wrong. Error: " + alert.message
        )
    })
})